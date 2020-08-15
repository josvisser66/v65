package asm

import "errors"

// Here is a good reference of 6502 addressing modes:
// http://www.obelisk.me.uk/6502/addressing.html

const (
	// Note: Do not change the values in this table, because they match
	// indices in slices in the opcodes maps.
	implicit        = 0  // <epsilon>
	accumulator     = 1  // A
	immediate       = 4  // #<expression>
	zeroPage        = 3  // <expression>
	zeroPageX       = 9  // <expression>, X
	zeroPageY       = 10 // <expression>, Y
	relative        = 11 // for branches
	absolute        = 2  // <expression>
	absoluteX       = 5  // <expression>, X
	absoluteY       = 6  // <expression>, Y
	indirect        = 12 // (<expression>)
	indexedIndirect = 7  // (<expression>, X)
	indirectIndexed = 8  // (<expression>), Y
)

var errorAddressingMode = errors.New("illegal addressing mode")

// parseIndirect parses addressing modes that start with an lparen. When
// the function is parsed the tokLparen needs to have already been parsed.
func (ctx *context) parseIndirect() (int, *exprValue, error) {
	val := ctx.expr()
	next := ctx.lexer.getToken()
	switch next.(type) {
	case *tokComma:
		// Must be (expr, X)
		next = ctx.lexer.getToken()
		if _, ok := next.(*tokRegisterX); !ok {
			ctx.error("expected X, not: '%T'", next)
			return -1, nil, errorAddressingMode
		}
		next = ctx.lexer.getToken()
		if _, ok := next.(*tokRightParen); !ok {
			ctx.error("expected ), not: '%T'", next)
			return -1, nil, errorAddressingMode
		}
		return indexedIndirect, val, nil
	case *tokRightParen:
		// Can be (expr) or (expr), Y
		next = ctx.lexer.getToken()
		if _, ok := next.(*tokNewLine); ok {
			// Is (expr) => Indirect.
			return indirect, val, nil
		}
		// Must be (expr), Y => Indirect indexed.
		if _, ok := next.(*tokComma); !ok {
			ctx.error("expected comma, not: '%T'", next)
			return -1, nil, errorAddressingMode
		}
		next = ctx.lexer.getToken()
		if _, ok := next.(*tokRegisterY); !ok {
			ctx.error("expected Y, not: '%T'", next)
			return -1, nil, errorAddressingMode
		}
		return indirectIndexed, val, nil
	}
	ctx.error("unexpected token: '%T'", next)
	return -1, nil, errorAddressingMode
}

// parseAbsolute parses addressing modes that start with an expression.
// Either expr or expr,X or expr,Y.
func (ctx *context) parseAbsolute() (int, *exprValue, error) {
	val := ctx.expr()
	next := ctx.lexer.getToken()
	switch next.(type) {
	case *tokNewLine:
		// If it was meant to be relative or zero page this gets resolved
		// later.
		return absolute, val, nil
	case *tokComma:
		next = ctx.lexer.getToken()
		switch next.(type) {
		case *tokRegisterX:
			return absoluteX, val, nil
		case *tokRegisterY:
			return absoluteY, val, nil
		}
	}
	ctx.error("unexpected token: '%T'", next)
	return -1, nil, errorAddressingMode
}

// parseAddressingMode parses an addressing mode. It does not give
// any errors about sizes or illegal use of external symbols. This
// will be done during stuffing the bytes into the code segment.
func (ctx *context) parseAddressingMode() (int, *exprValue, error) {
	tok := ctx.lexer.getToken()
	switch tok.(type) {
	case *tokNewLine:
		return implicit, nil, nil
	case *tokRegisterA:
		return accumulator, nil, nil
	case *tokHash:
		// Immediate addressing.
		return immediate, ctx.expr(), nil
	case *tokLeftParen:
		// Some form of indirect addressing.
		return ctx.parseIndirect()
	}
	// At this point we have either expr or expr,X or expr,Y.
	ctx.lexer.pushback(tok)
	return ctx.parseAbsolute()
}
