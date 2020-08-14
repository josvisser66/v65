package asm

import "testing"

func TestEqu(t *testing.T) {
	for _, tc := range []struct {
		str          string
		wantErrors   int
		wantWarnings int
		wantID       string
		wantValue    int64
	}{
		{"aap equ 42", 0, 0, "aap", 42},
		{"foo equ 42", 1, 0, "foo", 42},
		{"aap equ bar+1", 0, 0, "aap", 43},
		{"aap equ 1,1", 1, 0, "aap", 1},
		{"aap equ foo", 1, 0, "aap", 0},
		{"aap equ foo+1", 1, 0, "aap", 1},
		{"equ 5", 0, 1, "", 0},
		{"label equ *", 0, 0, "label", 0},
	} {
		println(tc.str)
		ctx := &context{
			lexer: &lexer{newSourceFromString(tc.str), nil},
			seg: newSegment(),
		}
		ctx.seg.symbols["foo"] = &externSymbol{"foo"}
		ctx.seg.symbols["bar"] = &localSymbol{"bar", 42, false}
		ctx.seg.symbols["baz"] = &localSymbol{"bar", 1000, false}
		ctx.assemble()
		if ctx.errors != tc.wantErrors {
			t.Errorf("assemble() errors; got:%d, want:%d", ctx.errors, tc.wantErrors)
		}
		if ctx.warnings != tc.wantWarnings {
			t.Errorf("assemble() errors; got:%d, want:%d", ctx.errors, tc.wantErrors)
		}
		if tc.wantID == "" {
			continue
		}
		sym, ok := ctx.seg.symbols[tc.wantID]
		if !ok {
			t.Errorf("symbol %s exists; got:false, want:true", tc.wantID)
		}
		ls, ok := sym.(*localSymbol)
		if !ok {
			t.Errorf("symbol %s is a local symbol; got:false, want:true", tc.wantID)
		}
		if ls.value != tc.wantValue {
			t.Errorf("symbol %s value; got:%d, want:%d", tc.wantID, ls.value, tc.wantValue)
		}
	}
}
