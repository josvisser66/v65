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

func TestOpcode(t *testing.T) {
	src := newSourceFromString("BRK")
	tok:= src.getToken()
	if op, ok := tok.(*tokOpcode); !ok {
		t.Errorf("token; got:'%T, want:'%T'", tok, &tokOpcode{})
	} else if op.opcode != "brk" {
		t.Errorf("op.opcde; got:'%s, want:'%s'", op.opcode, "brk")
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

func TestRune(t *testing.T) {
	src := newSourceFromString("'a'")
	tok := src.getToken()
	if tt, ok := tok.(*tokRune); !ok {
		t.Errorf("getToken(); got:%T, want:%T", tok, &tokRune{})
	} else if tt.r != 'a' {
		t.Errorf("tt; got:%c, want:%c", tt.r, 'a')
	}
}

func TestBadRune(t *testing.T) {
	bad := []string{"'b", "'cc'"}
	for _, h := range bad {
		src := newSourceFromString(h)
		tok := src.getToken()
		if _, ok := tok.(*tokError); !ok {
			t.Fatalf("getToken(); got:%T, want:%T", tok, &tokError{})
		}
	}
}

func TestString(t *testing.T) {
	src := newSourceFromString("\"abc\"")
	tok := src.getToken()
	if tt, ok := tok.(*tokString); !ok {
		t.Errorf("getToken(); got:%T, want:%T", tok, &tokString{})
	} else {
		if tt.s != "abc" {
			t.Errorf("tt.s; got:%s, want:%s", tt.s, "abc")
		}
	}
}

func TestStringEmbeddedQuote(t *testing.T) {
	src := newSourceFromString("\"a\"\"b\"\"c\"")
	tok := src.getToken()
	if tt, ok := tok.(*tokString); !ok {
		t.Errorf("getToken(); got:%T, want:%T", tok, &tokString{})
	} else {
		if tt.s != "a\"b\"c" {
			t.Errorf("tt.s; got:%s, want:%s", tt.s, "a\"b\"c")
		}
	}
}

func TestUnterminatedString(t *testing.T) {
	src := newSourceFromString("\"a\"\"b\"\"c")
	tok := src.getToken()
	if _, ok := tok.(*tokError); !ok {
		t.Errorf("getToken(); got:%T, want:%T", tok, &tokError{})
	}
}
