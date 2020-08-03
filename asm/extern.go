package asm

type tokExtern struct{}

func (*tokExtern) assemble(seg *segment, src *source, _ token) *source {
	for {
		// First, we expect an identifier.
		tok, ok := src.expect(seg, func(t token) bool {
			_, ok := t.(*tokIdentifier)
			return ok
		}, "identifier")
		if !ok {
			src.skipRestOfLine()
			return src
		}
		id := tok.(*tokIdentifier)
		if seg.symbols.register(id.id, &externSymbol{id.id}) {
			seg.warning(src, "redefinition of symbol %s\n", id.id)
		}
		// Then we either get a comma and we go around again, or we
		// get a newline and then we're done.
		tok = src.getToken()
		switch t := tok.(type) {
		case *tokNewLine:
			return src
		case *tokComma:
			// pass
		default:
			seg.error(src, "expected comma or end of line, not '%T'", t)
			src.skipToEOLN()
		}
	}
}

func init() {
	metaMap["extern"] = &tokExtern{}
}
