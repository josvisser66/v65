package asm

import "testing"

func TestEmit(t *testing.T) {
	for _, tc := range []struct {
		n int64
		want byte
	}{
		{1, 1},
		{2, 2},
		{255, 255},
		{-1, 255},
		{-2, 254},
		{1000, 0xe8},
	} {
		seg := newSegment()
		seg.emit(tc.n)
		if seg.code[0] != tc.want {
			t.Errorf("emit(%d); got:%d, want:%d", tc.n, seg.code[0], tc.want)
		}
	}
}

func TestEmitWord(t *testing.T) {
	for _, tc := range []struct {
		n int64
		want1 byte
		want2 byte
	}{
		{1, 0, 1},
		{2, 0, 2},
		{255, 0, 255},
		{-1, 255, 255},
		{-2, 255, 254},
		{1000, 3, 0xe8},
	} {
		seg := newSegment()
		seg.emitWord(tc.n)
		if seg.code[0] != tc.want1 || seg.code[1] != tc.want2 {
			t.Errorf("emit(%d); got:%d:%d, want:%d:%d", tc.n, seg.code[0], seg.code[1], tc.want1, tc.want2)
		}
	}
}