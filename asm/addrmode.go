package asm

const (
	errorAddrMode       = iota
	implied     // <epsilon>
	accumulator // A
	absolute    // <expression>
	immediate   // #<expression>
	absoluteX   // <expression>, X
	absoluteY   // <expression>, Y
	indirectX   // (<expression>), X
	indirectY   // (<expression>), Y
	indirect    // (<expression>)
)

func (ctx *context) expectNewline() {
	tok := ctx.src.getToken()
	if _, ok := tok.(*tokNewLine); ok {
		return
	}
	ctx.error("expected newline, not '%T'", tok)
}

// parseAddressingMode parses an addressing mode.
func (ctx *context) parseAddressingMode() (int, *exprValue) {
	tok := ctx.src.getToken()
	switch tok.(type) {
	case *tokNewLine:
		return implied, nil
	case *tokRegisterA:
		ctx.expectNewline()
		return accumulator, nil
	case *tokHash:
		val, next := ctx.expr(nil)
		if _, ok := next.(*tokNewLine); !ok {
			ctx.error("expected end-of-line, not '%T'", tok)
		}
		return immediate, val
	default:
		_, lparen := tok.(*tokLeftParen)
		val, next := ctx.expr(tok)
		if _, ok := next.(*tokNewLine); ok {
			if lparen {
				return indirect, val
			} else {
				return absolute, val
			}
		}
		if _, ok := next.(*tokComma); ok {
			reg := ctx.src.getToken()
			switch reg.(type) {
			case *tokRegisterX:
				ctx.expectNewline()
				if lparen {
					return indirectX, val
				}
				return absoluteX, val
			case *tokRegisterY:
				ctx.expectNewline()
				if lparen {
					return indirectY, val
				}
				return absoluteY, val
			}
			ctx.error("expected X or Y, not: '%T'", reg)
			ctx.src.skipRestOfLine()
			return errorAddrMode, nil
		}
		ctx.error("unexpected token: '%T'", next)
		ctx.src.skipRestOfLine()
		return errorAddrMode, nil
	}
}
