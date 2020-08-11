package asm

type tokGlobal struct{}

// assemble assembles a global statement.
func (*tokGlobal) assemble(ctx *context, label *localSymbol) {
	for {
		tok := ctx.src.getToken()
		id, ok := tok.(*tokIdentifier)
		if !ok {
			ctx.error("expected identifier, not '%T'", tok)
			ctx.src.skipRestOfLine()
			return
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
		tok = ctx.src.getToken()
		switch tok.(type) {
		case *tokNewLine:
			return
		case *tokComma:
			continue
		default:
			ctx.error("expected identifier, not '%T'", tok)
			ctx.src.skipRestOfLine()
			return
		}
	}
}

func init() {
	metaMap["global"] = &tokGlobal{}
}


