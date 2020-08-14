package asm

type tokExtern struct{}

func (*tokExtern) assemble(ctx *context, _label *localSymbol) error {
	for {
		// First, we expect an identifier.
		next := ctx.lexer.getToken()
		id, ok := next.(*tokIdentifier)
		if !ok {
			ctx.error("expected identifier, not: '%T(%v)'", next, next)
			return parseError
		}
		if ctx.seg.symbols.register(id.id, &externSymbol{id.id}) {
			ctx.warning("redefinition of symbol %s\n", id.id)
		}
		// Then we either get a comma and we go around again, or we
		// get a newline and then we're done.
		next = ctx.lexer.getToken()
		switch t := next.(type) {
		case *tokNewLine:
			return nil
		case *tokComma:
			// pass
		default:
			ctx.error("expected comma or end of line, not '%T'", t)
			return parseError
		}
	}
}

func init() {
	metaMap["extern"] = &tokExtern{}
}
