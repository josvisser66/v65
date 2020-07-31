package asm

import "testing"

func TestTokenizeIdentifiers(t *testing.T) {
	src := newSourceFromString("  aAp _NOOT Mies123 _wim12_foo ; comment")
	want := []string{"aap", "_noot", "mies123", "_wim12_foo"}
	got := make([]string, 0, 4)
	for {
		tok, s, eof := src.getToken()
		if eof {
			t.Errorf("eof; got:true, want:false")
			break
		}
		if tok == tokNewLine {
			break
		}
		if tok != tokIdentifier {
			t.Errorf("getToken(); got:%d, want:%d", tok, tokIdentifier)
		}
		got = append(got, s)
	}
	tok, _, _ := src.getToken()
	if tok != tokEOF {
		t.Errorf("getToken(); got:%d, want:%d", tok, tokEOF)
	}
	if len(got) != len(want) {
		t.Errorf("len(got); got:%d, want:%d", len(got), len(want))
	}
	for i:=0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("got[%d]; got:%s, want:%s", i, got[i], want[i])
		}
	}
}
