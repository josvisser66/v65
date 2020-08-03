package asm

import (
	"testing"
)

func TestExtern(t *testing.T) {
	for _, tc := range []struct {
		src       string
		wantErr   bool
		nSymbols  int
		nWarnings int
		symbols   []string
	}{
		{"extern foo", false, 1, 0, []string{"foo"}},
		{"extern foo, bar", false, 2, 0, []string{"foo", "bar"}},
		{"extern foo, bar, foo", false, 2, 1, []string{"foo", "bar"}},
		{"extern", true, 0, 0, []string{}},
		{"extern foo,", true, 1, 0, []string{"foo"}},
		{"extern foo,4", true, 1, 0, []string{"foo"}},
	} {
		println(tc.wantErr, tc.src)
		src := newSourceFromString(tc.src)
		seg, err := assemble(src)
		if (err != nil && !tc.wantErr) || (err == nil && tc.wantErr) {
			t.Errorf("assemble() err; got:%v, want-err:%v", err, tc.wantErr)
		}
		if seg == nil {
			continue
		}
		if seg.warnings != tc.nWarnings {
			t.Errorf("seg.warnings; got:%d, want:%d", seg.warnings, tc.nWarnings)
		}
		if len(seg.symbols) != tc.nSymbols {
			t.Errorf("len(seg.symbols); got:%d, want:%d", len(seg.symbols), tc.nSymbols)
		}
		for _, id := range tc.symbols {
			if sym, ok := seg.symbols[id]; !ok {
				t.Errorf("seg.symbols[%s] ok; got:%v, want:%v", id, ok, !ok)
			} else if _, ok := sym.(*externSymbol); !ok {
				t.Errorf("sym type; got:%T, want:%T", sym, &externSymbol{})
			}
		}
	}
}
