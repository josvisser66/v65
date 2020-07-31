package asm

import "testing"

func TestTokenizeIdentifiers(t *testing.T) {
	src := newSourceFromString("  aAp\t_NOOT\rMies123 _wim12_foo ; comment")
	want := []string{"aap", "_noot", "mies123", "_wim12_foo"}
	got := make([]string, 0, 4)
	for {
		tok, s, _, eof := src.getToken()
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
	tok, _, _, _ := src.getToken()
	if tok != tokEOF {
		t.Errorf("getToken(); got:%d, want:%d", tok, tokEOF)
	}
	if len(got) != len(want) {
		t.Errorf("len(got); got:%d, want:%d", len(got), len(want))
	}
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("got[%d]; got:%s, want:%s", i, got[i], want[i])
		}
	}
}

func TestTokenizeNumbers(t *testing.T) {
	src := newSourceFromString("0 010 0b101 0x4Af 42 $42 \t\r")
	got := make([]int64, 0, 6)
	want := []int64{0, 8, 5, 0x4af, 42, 0x42}
	for {
		tok, _, i, eof := src.getToken()
		if tok == tokNewLine {
			break
		}
		if eof {
			t.Fatalf("eof; got:true, want:false")
		}
		if tok != tokNumber {
			t.Fatalf("tok; got:%d, want:%d", tok, tokNumber)
		}
		got = append(got, i)
	}
	if len(got) != len(want) {
		t.Errorf("len(got); got:%d, want:%d", len(got), len(want))
	}
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("got[%d]; got:%d, want:%d", i, got[i], want[i])
		}
	}
}

func TestTokenizeBadHexNumbers(t *testing.T) {
	bad := []string{"0x", "$", "0xj", "$g"}
	for _, h := range bad {
		src := newSourceFromString(h)
		tok, s, _, eof := src.getToken()
		if eof {
			t.Fatalf("eof; got:true, want:false")
		}
		if tok != tokError {
			t.Fatalf("tok(%s); got:%d, want:%d", s, tok, tokError)
		}
	}
}
