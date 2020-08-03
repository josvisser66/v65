package asm

type tokExtern struct {}

func (tok *tokExtern) assemble(seg *segment, src *source, _ token) *source {
	for {
		tok := src.getToken()
		if seg.lexError(tok) {
			src.skipToEOLN()
			return src
		}
		if id, ok := tok.(*tokIdentifier); ok {
			if seg.symbols.register(id.id, &externSymbol{id.id}) {
				seg.warning(src, "redefinition of symbol %s\n", id.id)
			}
		}
		tok = src.getToken()
		switch t := tok.(type) {
		case *tokNewLine:
			return src
		case *tokRune:
			if t.r == ',' {
			continue
			}
		default:
			seg.error(src, "unexpected token: %T", t)
		}
	}
}

func init() {
	metaMap["extern"] = &tokExtern{}
}

