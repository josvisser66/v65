package asm

type exprValue struct {
	sym *externSymbol // If this is the value of a relocatable expression.
	val int64
}

func (ctx *context) expr() (val *exprValue, next token) {
	// Gets the token that starts the expression.
	tok := ctx.src.getToken()

	// If this token is an identifier, there is a change that we have
	// an expression that requires relocation. Such an expression starts
	// with an external symbol and then a positive or negative constant
	// offset.
	if id, ok := tok.(*tokIdentifier); ok {
		var sym symbol
		if sym, ok = ctx.seg.symbols[id.id]; !ok {
			// The error will be generated down there somewhere.
			v, next := ctx.level1(tok)
			return &exprValue{nil, v}, next
		}
		if externSym, ok := sym.(*externSymbol); ok {
			// This is an external symbol. The rest of the expression
			// can be + or - something else.
			next := ctx.src.getToken()
			var v int64
			if _, ok := next.(*tokPlus); ok {
				v, next = ctx.level1(nil)
			} else if _, ok := next.(*tokMinus); ok {
				v, next = ctx.level1(nil)
			}
			return &exprValue{externSym, v}, next
		}
		// The identifier is a label with a known value. Fallthrough.
	}

	v, next := ctx.level1(tok)
	return &exprValue{nil, v}, next
}

func (ctx *context) level1(nextToken token) (val int64, next token) {
	val, next = ctx.level2(nextToken)
	var v int64
	for {
		if next == nil {
			next = ctx.src.getToken()
		}
		if _, ok := next.(*tokOr); ok {
			v, next = ctx.level2(nil)
			val = val | v
		} else if _, ok := next.(*tokAnd); ok {
			v, next = ctx.level2(nil)
			val = val & v
		} else {
			return val, next
		}
	}
}

func (ctx *context) level2(nextToken token) (val int64, next token) {
	val, next = ctx.level3(nextToken)
	var v int64
	for {
		if next == nil {
			next = ctx.src.getToken()
		}
		if _, ok := next.(*tokPlus); ok {
			v, next = ctx.level3(nil)
			val = val + v
		} else if _, ok := next.(*tokMinus); ok {
			v, next = ctx.level3(nil)
				val = val - v
		} else {
			return val, next
		}
	}
}

func (ctx *context) level3(nextToken token) (val int64, next token) {
	val, next = ctx.level4(nextToken)
	var v int64
	for {
		if next == nil {
			next = ctx.src.getToken()
		}
		if _, ok := next.(*tokMultiply); ok {
			v, next = ctx.level4(nil)
			val = val * v
		} else if _, ok := next.(*tokDivide); ok {
			v, next = ctx.level4(nil)
			if v == 0 {
				ctx.error("division by zero")
			} else {
				val = val / v
			}
		} else {
			return val, next
		}
	}
}

func (ctx *context) level4(nextToken token) (val int64, next token) {
	if nextToken == nil {
		nextToken = ctx.src.getToken()
	}
	if ctx.lexError(nextToken) {
		ctx.src.skipRestOfLine()
		return 0,nil
	}
	if num, ok := nextToken.(*tokIntNumber); ok {
		return num.n, ctx.src.getToken()
	}
	if _, ok := nextToken.(*tokLeftParen); ok {
		v, next := ctx.level1(nil)
		if _, ok := next.(*tokRightParen); !ok {
			ctx.error("expected ')', not '%T'", next)
			ctx.src.skipToEOLN()
			return 0, nil
		}
		return v , nil
	}
	if id, ok := nextToken.(*tokIdentifier); ok {
		// Label, must be locally defined.
		sym, ok := ctx.seg.symbols[id.id]
		if !ok {
			ctx.error("unknown label: %s", id.id)
			return 0, nil
		}
		if localSym, ok := sym.(*localSymbol); ok {
			return localSym.value, nil
		}
		ctx.error("illegal use in expression of external label: %s", id.id)
		return 0, nil
	}
	if _, ok := nextToken.(*tokPlus); ok {
		// Unary plus operator.
		v, next := ctx.level4(nil)
		return v, next
	}
	if _, ok := nextToken.(*tokMinus); ok {
		// Unary minus operator.
		v, next := ctx.level4(nil)
		return -v, next
	}
	ctx.error("invalid expression; unexpected token: '%T'", nextToken)
	ctx.src.skipRestOfLine()
	return 0, nil
}