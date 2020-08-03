package asm

type tokDb struct {}
type tokDw struct {}
type tokDd struct {}
type tokDs struct {}

func (db *tokDb) assemble(seg *segment, src *source, _ token) *source {
	return src
}

func init() {
	metaMap["db"] = &tokDb{}
	metaMap["dw"] = &tokDw{}
	metaMap["dd"] = &tokDd{}
	metaMap["ds"] = &tokDs{}
}