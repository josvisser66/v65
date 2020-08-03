package asm

import "testing"

func (s *source) mustGetToken(t *testing.T) token {
	tok := s.getToken()
	if tt, ok := tok.(*tokError); ok {
		t.Fatalf("unexpected error: %v", tt)
	}
	return tok
}

func TestTokenizeIdentifiers(t *testing.T) {
	src := newSourceFromString("  aAp\t_NOOT\rMies123 _wim12_foo ; comment")
	want := []string{"aap", "_noot", "mies123", "_wim12_foo"}
	got := make([]string, 0, 4)
loop:
	for {
		switch tok := src.mustGetToken(t).(type) {
		case *tokNewLine:
			break loop
		case *tokEOF:
			t.Error("missing newline")
			break loop
		case *tokIdentifier:
			got = append(got, tok.id)
		default:
			t.Errorf("getToken(); got:%T, want:%T", tok, &tokIdentifier{})
		}
	}
	tok := src.mustGetToken(t)
	if _, ok := tok.(*tokEOF); !ok {
		t.Errorf("getToken(); got:%T, want:%T", tok, &tokEOF{})
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
loop:
	for {
		switch tok := src.mustGetToken(t).(type) {
		case *tokNewLine:
			break loop
		case *tokIntNumber:
			got = append(got, tok.n)
		default:
			t.Errorf("getToken(); got:%T, want:%T", tok, &tokIntNumber{})
		}
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
		tok := src.getToken()
		if _, ok := tok.(*tokError); !ok {
			t.Fatalf("getToken(); got:%T, want:%T", tok, &tokError{})
		}
	}
}
