package asm

// segment contains the generated machine language and symbols.
type segment struct {
	code    []byte
	lc      int
	top     int
	symbols symbolMap
}

// newSegment creates a new segment that can hold 64K of code and data.
// 64K should really be enough for everyone :-)
func newSegment() *segment {
	return &segment{
		code:    make([]byte, 65536),
		symbols: make(symbolMap),
	}
}

// emit writes a byte of data to the segment.
func (seg *segment) emit(b int64) {
	seg.code[seg.lc] = byte(b)
	seg.lc++
	if seg.lc > seg.top {
		seg.top = seg.lc
	}
}

// emitWord writes a word of data to the segment, big endian.
func (seg *segment) emitWord(w int64) {
	seg.code[seg.lc] = byte(w >> 8)
	seg.code[seg.lc+1] = byte(w & 255)
	seg.lc += 2
	if seg.lc > seg.top {
		seg.top = seg.lc
	}
}
