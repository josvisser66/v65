package asm

import "testing"

func TestExpressionEval(t *testing.T) {
	seg := newSegment()
	for _, tc := range []struct {
		str      string
		wantNum  int64
		wantErrors int
		wantNext func(t token) bool
	}{
		{"42", 42, 0,func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"42,", 42,0, func(t token) bool {
			_, ok := t.(*tokComma)
			return ok
		}},{"(42),", 42,0, func(t token) bool {
			_, ok := t.(*tokComma)
			return ok
		}},
		{"(42", 0,1, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
	} {
		println(tc.str, "=>", tc.wantNum)
		src := newSourceFromString(tc.str)
		saveErrors := seg.errors
		val, next := src.expr(seg)
		if next == nil {
			next = src.getToken()
		}
		if val.val != tc.wantNum {
			t.Errorf("expr(%s); got:%d, want:%d", tc.str, val.val, tc.wantNum)
		}
		if !tc.wantNext(next) {
			t.Error("tc.wantNext(t); got:false, want:true")
		}
		if saveErrors + tc.wantErrors != seg.errors {
			t.Errorf("seg.errors; got:%d, want:%d", seg.errors, saveErrors + tc.wantErrors)
		}
	}
}
