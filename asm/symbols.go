package asm

// symbol is an interface that represents a locally defined or external
// symbol.
type symbol interface {
}

// symbolMap is a map of identifier names to symbols.
type symbolMap map[string]symbol

// localSymbol is a symbol that is defined in this source file.
type localSymbol struct {
	id string
	value int64
	global bool // Should this symbol be exported?
}

// externSymbol is a symbol that is defined in another segment
// and the linker should resolve any uses.
type externSymbol struct {
	id string
}

// register registers a new symbol in the symbol table. Returns false
// if the symbol already exists. If that happens it will be replaced
// with the new one.
func (m symbolMap) register(id string, sym symbol) bool {
	 _, ok := m[id]
	m[id] = sym
	return ok
}