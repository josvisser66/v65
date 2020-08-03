package asm

import "fmt"

// segment contains the generated machine language.
type segment struct {
	code []byte
	lc int
	symbols symbolMap
	errors int
	warnings int
}

// newSegment creates a new segment that can hold 64K of code and data.
// 64K should really be enough for everyone :-)
func newSegment() *segment {
	return &segment {
		code: make([]byte, 0, 65536),
		symbols: make(symbolMap),
	}
}

// error reports an error.
func (seg *segment) error(src *source, s string, args ...interface{}) {
	fmt.Printf("[%d:%d] error: %s\n", src.lineNo, src.curPos, fmt.Sprintf(s, args...))
	seg.errors++
}

// lexError deals with the possibility of an error coming back from the lexer. Returns
// true if there was a lexer error (the error will already have been handled).
func (seg *segment) lexError(tok token) bool {
	if t, ok := tok.(*tokError); ok {
		fmt.Printf("[%d:%d] error: %s\n", t.lineNo, t.linePos, t.s)
		seg.errors++
		return true
	}
	return false
}

// warning reports a warning.
func (seg *segment) warning(src *source, s string, args ...interface{}) {
	fmt.Printf("[%d:%d] warning: %s", src.lineNo, src.curPos, fmt.Sprintf(s, args...))
	seg.warnings++
}

// emit writes a byte of data to the segment.
func (seg *segment) emit(b byte) {
	seg.code[seg.lc] = b
	seg.lc++
}

// emitWord writes a word of data to the segment, big endian.
func (seg *segment) emitWord(w uint16) {
	seg.code[seg.lc] = byte(w >> 8)
	seg.code[seg.lc+1] = byte(w & 255)
	seg.lc += 2
}