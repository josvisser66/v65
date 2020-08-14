package asm

import (
	"testing"
)

func TestExtern(t *testing.T) {
	for _, tc := range []struct {
		src       string
		wantErrors int
		wantWarnings int
		wantSymbols  int
		symbols   []string
	}{
		{"extern foo", 0, 0, 1, []string{"foo"}},
		{"extern foo, bar", 0, 0, 2, []string{"foo", "bar"}},
		{"extern foo, bar, foo", 0, 1, 2, []string{"foo", "bar"}},
		{"extern", 1, 0, 0, []string{}},
		{"extern foo,", 1, 0, 1, []string{"foo"}},
		{"extern foo,4,bar", 1, 0, 1, []string{"foo"}},
	} {
		println(tc.src)
		ctx := &context{
			lexer: &lexer{newSourceFromString(tc.src), nil},
			seg: newSegment(),
		}
		ctx.assemble()
		if ctx.errors != tc.wantErrors {
			t.Errorf("assemble() errors; got:%d, want:%d", ctx.errors, tc.wantErrors)
		}
		if ctx.warnings != tc.wantWarnings {
			t.Errorf("assemble() warnings; got:%d, want:%d", ctx.warnings, tc.wantWarnings)
		}
		if len(ctx.seg.symbols) != tc.wantSymbols {
			t.Errorf("len(seg.symbols); got:%d, want:%d", len(ctx.seg.symbols), tc.wantSymbols)
		}
		for _, id := range tc.symbols {
			if sym, ok := ctx.seg.symbols[id]; !ok {
				t.Errorf("seg.symbols[%s] ok; got:%v, want:%v", id, ok, !ok)
			} else if _, ok := sym.(*externSymbol); !ok {
				t.Errorf("sym type; got:%T, want:%T", sym, &externSymbol{})
			}
		}
	}
}
