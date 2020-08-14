package asm

// Here is a good reference of 6502 addressing modes:
// http://www.obelisk.me.uk/6502/addressing.html

const (
	errorAddrMode   = iota
	implicit        // <epsilon>
	accumulator     // A
	immediate       // #<expression>
	zeroPage        // <expression>
	zeroPageX       // <expression>, X
	zeroPageY       // <expression>, Y
	relative        // for branches
	absolute        // <expression>
	absoluteX       // <expression>, X
	absoluteY       // <expression>, Y
	indirect        // (<expression>)
	indexedIndirect // (<expression>, X)
	indirectIndexed // (<expression>), Y
)

func (ctx *context) expectNewline(tok token) {
	tok = ctx.src.getNextToken(tok)
	if _, ok := tok.(*tokNewLine); ok {
		return
	}
	ctx.error("expected end of line, not '%T'", tok)
}

// parseIndirect parses addressing modes that start with an lparen. When
// the function is parsed the tokLparen needs to have already been parsed.
func (ctx *context) parseIndirect() (int, *exprValue) {
	val, next := ctx.expr(nil)
	next = ctx.src.getNextToken(next)
	switch next.(type) {
	case *tokComma:
		// Must be (expr, X)
		next = ctx.src.getToken()
		if _, ok := next.(*tokRegisterX); !ok {
			ctx.error("expected X, not: '%T'", next)
			ctx.src.skipRestOfLine()
			return errorAddrMode, nil
		}
		next = ctx.src.getToken()
		if _, ok := next.(*tokRightParen); !ok {
			ctx.error("expected ), not: '%T'", next)
			ctx.src.skipRestOfLine()
			return errorAddrMode, nil
		}
		ctx.expectNewline(nil)
		return indexedIndirect, val
	case *tokRightParen:
		// Can be (expr) or (expr), Y
		next = ctx.src.getToken()
		if _, ok := next.(*tokNewLine); ok {
			// Is (expr) => Indirect.
			return indirect, val
		}
		// Must be (expr), Y => Indirect indexed.
		if _, ok := next.(*tokComma); !ok {
			ctx.error("expected comma, not: '%T'", next)
			ctx.src.skipRestOfLine()
			return errorAddrMode, nil
		}
		next = ctx.src.getToken()
		if _, ok := next.(*tokRegisterY); !ok {
			ctx.error("expected Y, not: '%T'", next)
			ctx.src.skipRestOfLine()
			return errorAddrMode, nil
		}
		ctx.expectNewline(nil)
		return indirectIndexed, val
	}
	ctx.error("unexpected token: '%T'", next)
	ctx.src.skipRestOfLine()
	return errorAddrMode, nil
}

// parseAbsolute parses addressing modes that start with an expression.
// Either expr or expr,X or expr,Y.
func (ctx *context) parseAbsolute(firstToken token) (int, *exprValue) {
	val, next := ctx.expr(firstToken)
	next = ctx.src.getNextToken(next)
	switch next.(type) {
	case *tokNewLine:
		// If it was meant to be relative or zero page this gets resolved
		// later.
		return absolute, val
	case *tokComma:
		next = ctx.src.getToken()
		switch next.(type) {
		case *tokRegisterX:
			return absoluteX, val
		case *tokRegisterY:
			return absoluteY, val
		}
	}
	ctx.error("unexpected token: '%T'", next)
	return errorAddrMode, nil
}

// parseAddressingMode parses an addressing mode.
func (ctx *context) parseAddressingMode() (int, *exprValue) {
	tok := ctx.src.getToken()
	switch tok.(
	type
	) {
	case *tokNewLine:
		return implicit, nil
	case *tokRegisterA:
		ctx.expectNewline(nil)
		return accumulator, nil
	case *tokHash:
		// Immediate addressing.
		val, next := ctx.expr(nil)
		ctx.expectNewline(next)
		return immediate, val
	case *tokLeftParen:
		// Some form of indirect addressing.
		return ctx.parseIndirect()
	}
	// At this point we have either expr or expr,X or expr,Y.
	return ctx.parseAbsolute(tok)
}
