package asm

// segment contains the generated machine language and symbols.
type segment struct {
	code []byte
	lc int
	symbols symbolMap
}

// newSegment creates a new segment that can hold 64K of code and data.
// 64K should really be enough for everyone :-)
func newSegment() *segment {
	return &segment {
		code: make([]byte, 0, 65536),
		symbols: make(symbolMap),
	}
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