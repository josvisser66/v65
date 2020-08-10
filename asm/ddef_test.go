package asm

import "testing"

func TestDb(t *testing.T) {
	for _, tc := range []struct {
		str        string
		wantErrors int
		wantRelocs int
		relocSize int
		wantBytes  []byte
	}{
		{"db 1,2,-1", 0, 0, 1,[]byte{1, 2, 255}},
		{"db 1,2,", 1, 0,1, []byte{1, 2, 0}},
		{"db foo, bar, baz", 0, 1,1, []byte{0, 42, 0xe8}},
		{"db foo+2", 0, 1, 1, []byte{2}},
		{"dw 1000", 0, 0, 2, []byte{0x3, 0xe8}},
		{"dw foo+1000", 0, 1, 2, []byte{0x3, 0xe8}},
		{"dd 65538", 0, 0, 4, []byte{0, 1, 0, 2}},
		{"dd foo+2", 0, 1, 4, []byte{0, 0, 0, 2}},
		{"dd 0x12345678", 0, 0,4,  []byte{0x12, 0x34, 0x56, 0x78}},
		{"dd 0x87654321", 0, 0, 4, []byte{0x87, 0x65, 0x43, 0x21}},
		{"ds \"abc\",\"def\"", 0, 0, 0, []byte{97, 98, 99, 100, 101, 102}},
		{"ds \"abc\",", 1, 0, 0, []byte{97, 98, 99}},
		{"db", 1, 0, 1, []byte{0}},
	} {
		println(tc.str)
		ctx := &context{
			src: newSourceFromString(tc.str),
			seg: newSegment(),
		}
		ctx.seg.symbols["foo"] = &externSymbol{"foo"}
		ctx.seg.symbols["bar"] = &localSymbol{"bar", 42, 1, false}
		ctx.seg.symbols["baz"] = &localSymbol{"bar", 1000, 2, false}
		ctx.assemble()
		if ctx.errors != tc.wantErrors {
			t.Errorf("assemble() errors; got:%d, want:%d", ctx.errors, tc.wantErrors)
		}
		if ctx.seg.size != len(tc.wantBytes) {
			t.Errorf("ctx.seg.size; got:%d, want:%d", ctx.seg.size, len(tc.wantBytes))
		}
		if len(ctx.seg.relocs) != tc.wantRelocs {
			t.Errorf("len(ctx.seg.relocs); got:%d, want:%d", len(ctx.seg.relocs), tc.wantRelocs)
		}
		for _, relocs := range ctx.seg.relocs {
			for _, reloc := range relocs {
				if reloc.size != tc.relocSize {
					t.Errorf("reloc.size; got:%d, want:%d", reloc.size, 1)
				}
			}
		}
		for i, b := range tc.wantBytes {
			if ctx.seg.code[i] != b {
				t.Errorf("code[%d]; got:%d, want:%d", i, ctx.seg.code[i], b)
			}
		}
	}
}
