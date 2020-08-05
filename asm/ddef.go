package asm

type tokDb struct{}
type tokDw struct{}
type tokDd struct{}
type tokDs struct{}

func assembleDdef(ctx *context, size int, emit func(int64)) {
	for {
		val, next := ctx.expr()
		ctx.seg.relocs.maybeAdd(val, ctx.seg.lc, size)
		emit(val.val)

		if _, ok := next.(*tokNewLine); ok {
			return
		}
		if _, ok := next.(*tokComma); ok {
			continue
		}
		ctx.error("expected ',' or newline, got: '%T'", next)
		ctx.src.skipRestOfLine()
		return
	}
}

func (*tokDb) assemble(ctx *context) {
	assembleDdef(ctx, 1, func(n int64) { ctx.seg.emit(n) })
}

func (*tokDw) assemble(ctx *context) {
	assembleDdef(ctx, 2, func(n int64) { ctx.seg.emitWord(n) })
}

func (*tokDd) assemble(ctx *context) {
	assembleDdef(ctx, 4, func(n int64) { ctx.seg.emitDWord(n) })
}

func init() {
	metaMap["db"] = &tokDb{}
	metaMap["dw"] = &tokDw{}
	metaMap["dd"] = &tokDd{}
	metaMap["ds"] = &tokDs{}
}
