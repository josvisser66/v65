package asm

type relocation struct {
	lc int
	size int
}

type relocMap map[string][]relocation

func (r relocMap) add(sym string, lc int, size int) {
	_, ok := r[sym]
	if !ok {
		r[sym] = make([]relocation, 0, 1)
	}
	r[sym] = append(r[sym], relocation{lc, size})
}

func (r relocMap) maybeAdd(val *exprValue, lc int, size int) {
	if val.sym != nil {
		r.add(val.sym.id, lc, size)
	}
}