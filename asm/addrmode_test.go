package asm

import "testing"

func TestAdrMode(t *testing.T) {
	var f42  int64 = 42
	var f43  int64 = 43
	for _, tc := range []struct {
		str          string
		wantError bool
		wantMode   int
		wantValue    *int64
	}{
		{"", false, implicit, nil},
		{"A",false, accumulator, nil},
		{"#42",false, immediate, &f42},
		{"#42,",true, immediate, &f42},
		{"43",false, absolute, &f43},
		{"43a",true, errorAddrMode, nil},
		{"43,X",false, absoluteX, &f43},
		{"43,Y",false, absoluteY, &f43},
		{"43,Z",true, errorAddrMode, nil},
		{"(42)",false, indirect, &f42},
		{"(42):z",true, errorAddrMode, nil},
		{"(42, X)",false, indexedIndirect, &f42},
		{"(42),Y",false, indirectIndexed, &f42},
		{"(42),Z",true, errorAddrMode, nil},
	} {
		println(tc.str)
		ctx := &context{
			lexer: &lexer{newSourceFromString(tc.str), nil},
			seg: newSegment(),
		}
		ctx.seg.symbols["foo"] = &externSymbol{"foo"}
		ctx.seg.symbols["bar"] = &localSymbol{"bar", 42, false}
		ctx.seg.symbols["baz"] = &localSymbol{"bar", 1000, false}
		mode, val := ctx.parseAddressingMode()
		if ctx.errors != 0 && ! tc.wantError {
			t.Errorf("parseAddressingMode() errors; got:%d, want:0", ctx.errors)
		}
		if mode != tc.wantMode {
			t.Errorf("parseAddressingMode() mode; got:%d, want:%d", mode, tc.wantMode)
		}
		if val == nil && tc.wantValue != nil {
			t.Error("parseAddressingMode() value; got:nil, want:non-nil")
			continue
		}
		if val != nil && tc.wantValue == nil {
			t.Error("parseAddressingMode() value; got:non-nil, want:nil")
			continue
		}
		if val != nil && val.val != *tc.wantValue {
			t.Errorf("parseAddressingMode() value; got:%d, want:%d", val.val, tc.wantValue)
		}
	}
}

