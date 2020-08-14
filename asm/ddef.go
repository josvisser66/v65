package asm

type tokDb struct{}
type tokDw struct{}
type tokDd struct{}
type tokDs struct{}

// assembleDDef assembles db, dw, and dd instructions.
func assembleDdef(ctx *context, size int, emit func(int64)) error {
	for {
		val := ctx.expr()
		next := ctx.lexer.getToken()
		ctx.seg.relocs.maybeAdd(val, ctx.seg.lc, size)
		emit(val.val)
		if _, ok := next.(*tokNewLine); ok {
			return nil
		}
		if _, ok := next.(*tokComma); ok {
			continue
		}
		ctx.error("expected ',' or newline, got: '%T'", next)
		ctx.lexer.src.skipToEOLN()
		return parseError
	}
}

// assemble assembles a db instruction.
func (*tokDb) assemble(ctx *context, _label *localSymbol) error {
	return assembleDdef(ctx, 1, func(n int64) { ctx.seg.emit(n) })
}

// assemble assembles a dw instruction.
func (*tokDw) assemble(ctx *context, _label *localSymbol) error {
	return assembleDdef(ctx, 2, func(n int64) { ctx.seg.emitWord(n) })
}

// assemble assembles a dd instruction.
func (*tokDd) assemble(ctx *context, _label *localSymbol) error {
	return assembleDdef(ctx, 4, func(n int64) { ctx.seg.emitDWord(n) })
}

// assemble assembles a ds instruction.
func (*tokDs) assemble(ctx *context, _label *localSymbol) error {
	for {
		tok := ctx.lexer.getToken()
		tt, ok := tok.(*tokString)
		if !ok {
			ctx.error("expected string")
			return parseError
		}
		for _, b := range []byte(tt.s) {
			ctx.seg.emit(int64(b))
		}
		next := ctx.lexer.getToken()
		if _, ok := next.(*tokNewLine); ok {
			return nil
		}
		if _, ok := next.(*tokComma); ok {
			continue
		}
		ctx.error("expected ',' or newline, got: '%T'", next)
		ctx.lexer.src.skipToEOLN()
		return parseError
	}
}

func init() {
	metaMap["db"] = &tokDb{}
	metaMap["dw"] = &tokDw{}
	metaMap["dd"] = &tokDd{}
	metaMap["ds"] = &tokDs{}
}
