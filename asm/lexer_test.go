package asm

import "testing"

func (l *lexer) mustGetToken(t *testing.T) token {
	tok := l.getToken()
	if tt, ok := tok.(*tokError); ok {
		t.Fatalf("unexpected lexer error: %v", tt)
	}
	return tok
}

func (l *lexer) mustReadNewlines(t *testing.T, n int) {
	for i :=0; i < n; i++ {
		tok := l.getToken()
		if _, ok := tok.(*tokNewLine); !ok {
			t.Fatalf("tok; got:%T, want:%T", tok, &tokNewLine{})
		}
	}
}

func TestNewlineProcessing(t *testing.T) {
	lexer := &lexer{src : newSourceFromString("\n\naap")}
	lexer.mustReadNewlines(t, 1000)
	lexer.src.moveToNextLine()
	lexer.mustReadNewlines(t, 1000)
	lexer.src.moveToNextLine()
	tok := lexer.mustGetToken(t)
	if tt, ok := tok.(*tokIdentifier); !ok {
		t.Errorf("tok; got:%T, want:%T", tok, &tokIdentifier{})
	} else if tt.id != "aap" {
		t.Errorf("id; got:%s, want:aap", tt.id)
	}
}

func TestTokenizeIdentifiers(t *testing.T) {
	lexer := &lexer{src : newSourceFromString("  aAp\t_NOOT\rMies123 _wim12_foo ; comment")}
	want := []string{"aap", "_noot", "mies123", "_wim12_foo"}
	got := make([]string, 0, 4)
loop:
	for {
		switch tok := lexer.mustGetToken(t).(type) {
		case *tokNewLine:
			break loop
		case *tokEOF:
			t.Error("missing newline")
			break loop
		case *tokIdentifier:
			got = append(got, tok.id)
		default:
			t.Errorf("getToken(); got:%T(%v), want:%T", tok, tok, &tokIdentifier{})
		}
	}
	lexer.src.moveToNextLine()
	tok := lexer.mustGetToken(t)
	if _, ok := tok.(*tokEOF); !ok {
		t.Errorf("getToken(); got:%T(%v), want:%T", tok, tok, &tokEOF{})
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
	lexer := &lexer{src : newSourceFromString("BRK")}
	tok:= lexer.getToken()
	if op, ok := tok.(*tokOpcode); !ok {
		t.Errorf("token; got:'%T, want:'%T'", tok, &tokOpcode{})
	} else if op.opcode != "brk" {
		t.Errorf("op.opcde; got:'%s, want:'%s'", op.opcode, "brk")
	}
}

func TestTokenizeNumbers(t *testing.T) {
	lexer := &lexer{src : newSourceFromString("0 010 0b101 0x4Af 42 $42 \t\r")}
	got := make([]int64, 0, 6)
	want := []int64{0, 8, 5, 0x4af, 42, 0x42}
loop:
	for {
		switch tok := lexer.mustGetToken(t).(type) {
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
		lexer := &lexer{src : newSourceFromString(h)}
		tok := lexer.getToken()
		if _, ok := tok.(*tokError); !ok {
			t.Fatalf("getToken(); got:%T, want:%T", tok, &tokError{})
		}
	}
}

func TestRune(t *testing.T) {
	lexer := &lexer{src : newSourceFromString("'a'")}
	tok := lexer.getToken()
	if tt, ok := tok.(*tokRune); !ok {
		t.Errorf("getToken(); got:%T, want:%T", tok, &tokRune{})
	} else if tt.r != 'a' {
		t.Errorf("tt; got:%c, want:%c", tt.r, 'a')
	}
}

func TestBadRune(t *testing.T) {
	bad := []string{"'b", "'cc'"}
	for _, h := range bad {
		lexer := &lexer{src : newSourceFromString(h)}
		tok := lexer.getToken()
		if _, ok := tok.(*tokError); !ok {
			t.Fatalf("getToken(); got:%T, want:%T", tok, &tokError{})
		}
	}
}

func TestString(t *testing.T) {
	lexer := &lexer{src : newSourceFromString("\"abc\"")}
	tok := lexer.getToken()
	if tt, ok := tok.(*tokString); !ok {
		t.Errorf("getToken(); got:%T, want:%T", tok, &tokString{})
	} else {
		if tt.s != "abc" {
			t.Errorf("tt.s; got:%s, want:%s", tt.s, "abc")
		}
	}
}

func TestStringEmbeddedQuote(t *testing.T) {
	lexer := &lexer{src : newSourceFromString("\"a\"\"b\"\"c\"")}
	tok := lexer.getToken()
	if tt, ok := tok.(*tokString); !ok {
		t.Errorf("getToken(); got:%T, want:%T", tok, &tokString{})
	} else {
		if tt.s != "a\"b\"c" {
			t.Errorf("tt.s; got:%s, want:%s", tt.s, "a\"b\"c")
		}
	}
}

func TestUnterminatedString(t *testing.T) {
	lexer := &lexer{src : newSourceFromString("\"a\"\"b\"\"c")}
	tok := lexer.getToken()
	if _, ok := tok.(*tokError); !ok {
		t.Errorf("getToken(); got:%T, want:%T", tok, &tokError{})
	}
}
