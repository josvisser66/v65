package asm

type tokGlobal struct{}

// assemble assembles a global statement.
func (*tokGlobal) assemble(ctx *context, label *localSymbol) error {
	for {
		next := ctx.lexer.getToken()
		id, ok := next.(*tokIdentifier)
		if !ok {
			ctx.error("expected identifier, not '%T'", next)
			return parseError
		}
		sym, ok := ctx.seg.symbols[id.id]
		if !ok {
			ctx.error("undefined symbol %s", id.id)
		} else {
			ls, ok := sym.(*localSymbol)
			if !ok {
				ctx.error("cannot make an external symbol global")
			} else {
				ls.global = true
			}
		}
		next = ctx.lexer.getToken()
		switch next.(type) {
		case *tokNewLine:
			return nil
		case *tokComma:
			// pass
		default:
			ctx.error("expected identifier, not '%T'", next)
			return parseError
		}
	}
}

func init() {
	metaMap["global"] = &tokGlobal{}
}


