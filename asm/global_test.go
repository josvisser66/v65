package asm

import "testing"

func TestGlobal(t *testing.T) {
	for _, tc := range []struct {
		str        string
		wantErrors int
		wantIDs    []string
	}{
		{"global bar", 0, []string{"bar"}},
		{"global bar, baz", 0, []string{"bar", "baz"}},
		{"global foo", 1, []string{}},
		{"global foo,", 2, []string{}},
		{"global foo,4", 2, []string{}},
		{"global bar, foo, baz, foo", 2, []string{"bar", "baz"}},
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
		for _, wantID := range tc.wantIDs {
			sym, ok := ctx.seg.symbols[wantID]
			if !ok {
				t.Errorf("symbol %s exists; got:false, want:true", wantID)
			}
			ls, ok := sym.(*localSymbol)
			if !ok {
				t.Errorf("symbol %s is a local symbol; got:false, want:true", wantID)
			} else if !ls.global {
				t.Errorf("symbol %s is global; got:false, want:true", wantID)
			}
		}
	}
}
