package asm

type tokEqu struct{}

func (*tokEqu) assemble(ctx *context, label *localSymbol) {
	val, next := ctx.expr()
	if _, ok := next.(*tokNewLine); !ok {
		ctx.error("expected end-of-line")
		ctx.src.skipRestOfLine()
	}
	if val.sym!=nil {
		ctx.error("defining a local symbol with an external value is not allowed")
	}
	if label == nil {
		ctx.warning("equ without label, value is lost")
	} else {
		label.value = val.val
	}
}

func init() {
	metaMap["equ"] = &tokEqu{}
}

