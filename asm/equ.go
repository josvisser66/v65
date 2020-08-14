package asm

type tokEqu struct{}

func (*tokEqu) assemble(ctx *context, label *localSymbol) error {
	val := ctx.expr()
	next :=  ctx.lexer.getToken()
	if _, ok := next.(*tokNewLine); !ok {
		return nil
	}
	if val.sym!=nil {
		ctx.error("defining a local symbol with an external value is not allowed")
		return parseError
	}
	if label == nil {
		ctx.warning("equ without label, value is lost")
	} else {
		label.value = val.val
	}
	return nil
}

func init() {
	metaMap["equ"] = &tokEqu{}
}

