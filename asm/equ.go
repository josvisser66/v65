package asm

type tokEqu struct{}

func (*tokEqu) assemble(ctx *context, label *localSymbol) (err error) {
	val := ctx.expr()
	if val.sym != nil {
		ctx.error("defining a local symbol with an external value is not allowed")
		err = parseError
	}
	if label == nil {
		ctx.warning("equ without label, value is lost")
	} else {
		label.value = val.val
	}
	return err
}

func init() {
	metaMap["equ"] = &tokEqu{}
}
