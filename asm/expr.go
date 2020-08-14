package asm

type exprValue struct {
	sym *externSymbol // If this is the value of a relocatable expression.
	val int64
}

func (ctx *context) expr() *exprValue {
	tok := ctx.lexer.getToken()

	// If this token is an identifier, there is a change that we have
	// an expression that requires relocation. Such an expression starts
	// with an external symbol and then a positive or negative constant
	// offset.
	if id, ok := tok.(*tokIdentifier); ok {
		var sym symbol
		if sym, ok = ctx.seg.symbols[id.id]; !ok {
			// The error will be generated down there somewhere.
			ctx.lexer.pushback(tok)
			return &exprValue{nil, ctx.level1()}
		}
		if externSym, ok := sym.(*externSymbol); ok {
			// This is an external symbol. The rest of the expression
			// can be + or - something else.
			next := ctx.lexer.getToken()
			var v int64
			if _, ok := next.(*tokPlus); ok {
				v = ctx.level1()
			} else if _, ok := next.(*tokMinus); ok {
				v = -ctx.level1()
			} else {
				ctx.lexer.pushback(next)
			}
			return &exprValue{externSym, v}
		}
		// The identifier is a label with a known value. Fallthrough.
	}

	ctx.lexer.pushback(tok)
	v := ctx.level1()
	return &exprValue{nil, v}
}

func (ctx *context) level1() int64 {
	val := ctx.level2()
	for {
		next := ctx.lexer.getToken()
		if _, ok := next.(*tokOr); ok {
			val = val | ctx.level2()
		} else if _, ok := next.(*tokAnd); ok {
			val = val & ctx.level2()
		} else {
			ctx.lexer.pushback(next)
			return val
		}
	}
}

func (ctx *context) level2() int64 {
	val := ctx.level3()
	for {
		next := ctx.lexer.getToken()
		if _, ok := next.(*tokPlus); ok {
			val = val + ctx.level3()
		} else if _, ok := next.(*tokMinus); ok {
			val = val - ctx.level3()
		} else {
			ctx.lexer.pushback(next)
			return val
		}
	}
}

func (ctx *context) level3() int64 {
	val := ctx.level4()
	for {
		next := ctx.lexer.getToken()
		if _, ok := next.(*tokMultiply); ok {
			val = val * ctx.level4()
		} else if _, ok := next.(*tokDivide); ok {
			v := ctx.level4()
			if v == 0 {
				ctx.error("division by zero")
			} else {
				val = val / v
			}
		} else {
			ctx.lexer.pushback(next)
			return val
		}
	}
}

func (ctx *context) level4() int64 {
	next := ctx.lexer.getToken()
	if ctx.lexError(next) {
		return 0
	}
	if _,  ok := next.(*tokMultiply); ok {
		// Current location counter.
		return int64(ctx.seg.lc)
	}
	if num, ok := next.(*tokIntNumber); ok {
		return num.n
	}
	if _, ok := next.(*tokLeftParen); ok {
		v := ctx.level1()
		next := ctx.lexer.getToken()
		if _, ok := next.(*tokRightParen); !ok {
			ctx.error("expected ')', not '%T'", next)
			return 0
		}
		return v
	}
	if id, ok := next.(*tokIdentifier); ok {
		// Label, must be locally defined.
		sym, ok := ctx.seg.symbols[id.id]
		if !ok && ctx.pass == 2 {
			ctx.error("unknown label: %s", id.id)
			return 0
		}
		if localSym, ok := sym.(*localSymbol); ok {
			return localSym.value
		}
		ctx.error("illegal use in expression of external label: %s", id.id)
		return 0
	}
	if _, ok := next.(*tokPlus); ok {
		// Unary plus operator.
		return ctx.level4()
	}
	if _, ok := next.(*tokMinus); ok {
		// Unary minus operator.
		return -ctx.level4()
	}
	ctx.lexer.pushback(next)
	ctx.error("invalid expression; unexpected token: '%T(%v)'", next, next)
	return 0
}