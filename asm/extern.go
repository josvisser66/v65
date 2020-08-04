package asm

type tokExtern struct{}

func (*tokExtern) assemble(ctx *context) {
	for {
		// First, we expect an identifier.
		tok, ok := ctx.expect(func(t token) bool {
			_, ok := t.(*tokIdentifier)
			return ok
		}, "identifier")
		if !ok {
			ctx.src.skipRestOfLine()
			return
		}
		id := tok.(*tokIdentifier)
		if ctx.seg.symbols.register(id.id, &externSymbol{id.id}) {
			ctx.warning("redefinition of symbol %s\n", id.id)
		}
		// Then we either get a comma and we go around again, or we
		// get a newline and then we're done.
		tok = ctx.src.getToken()
		switch t := tok.(type) {
		case *tokNewLine:
			return
		case *tokComma:
			// pass
		default:
			ctx.error("expected comma or end of line, not '%T'", t)
			ctx.src.skipToEOLN()
		}
	}
}

func init() {
	metaMap["extern"] = &tokExtern{}
}
