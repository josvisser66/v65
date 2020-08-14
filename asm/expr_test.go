package asm

import "testing"

func TestExpressionEval(t *testing.T) {
	seg := newSegment()
	seg.symbols["fu"] = &externSymbol{"fu"}
	seg.symbols["bar"] = &localSymbol{"bar", 7, false}
	for _, tc := range []struct {
		str        string
		wantNum    int64
		wantSym bool
		wantErrors int
		wantNext   func(t token) bool
	}{
		{"42", 42, false,0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"-40*2/4", -20, false,0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"-40*2/4+10*7", 50, false,0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"65535 & 255", 255, false,0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"2 | 1", 3, false,0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"+42", 42, false,0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"-42", -42, false,0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"--42", 42, false,0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"fu", 0, true, 0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"bar", 7, false, 0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"baz", 0, false, 1, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"42", 42, false, 0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"*+5", 5, false, 0, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{"42,", 42, false, 0, func(t token) bool {
			_, ok := t.(*tokComma)
			return ok
		}}, {"(42),", 42, false, 0, func(t token) bool {
			_, ok := t.(*tokComma)
			return ok
		}},
		{"(42", 0, false, 1, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
		{",", 0, false, 1, func(t token) bool {
			_, ok := t.(*tokComma)
			return ok
		}},
		{"", 0, false, 1, func(t token) bool {
			_, ok := t.(*tokNewLine)
			return ok
		}},
	} {
		println(tc.str, "=>", tc.wantNum)
		ctx := &context{
			lexer: &lexer{newSourceFromString(tc.str), nil},
			seg: seg,
		}
		val := ctx.expr()
		next := ctx.lexer.getToken()
		if (val.sym == nil && tc.wantSym) || (val.sym!= nil && !tc.wantSym) {
			t.Errorf("val.sym: got:%v, want-nil:%v", val.sym, tc.wantSym)
		}
		if val.val != tc.wantNum {
			t.Errorf("expr(%s); got:%d, want:%d", tc.str, val.val, tc.wantNum)
		}
		if !tc.wantNext(next) {
			t.Error("tc.wantNext(t); got:false, want:true")
		}
		if tc.wantErrors != ctx.errors {
			t.Errorf("seg.errors; got:%d, want:%d", ctx.errors, tc.wantErrors)
		}
	}
}
