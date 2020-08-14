package asm

import "testing"

const testData = "aap\n\nnoot"

func testIt(t *testing.T, src *source) {
	got := make([]rune, 0, 10)
	for {
		r, eof := src.consumeRune()
		if eof {
			break
		}
		got = append(got, r)
		if r == '\n' {
			src.moveToNextLine()
		}
	}
	want := []rune{'a', 'a', 'p', '\n', '\n', 'n', 'o', 'o', 't', '\n'}
	if len(got) != len(want) {
		t.Errorf("len(got); got:%d, want:%d", len(got), len(want))
	}
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("got[%d]; got:%d, want:%d", i, got[i], want[i])
		}
	}
}

func TestConsumeRune(t *testing.T) {
	src := newSourceFromString(testData)
	testIt(t, src)
}

func TestPeekRune(t *testing.T) {
	src := newSourceFromString(testData)
	for i := 0; i < 3; i++ {
		r, eof := src.peekRune()
		if r != 'a' || eof {
			t.Fatalf("peekRune(); got:%d,%v, want:'a',false", r, eof)
		}
	}
	testIt(t, src)
	for i := 0; i < 3; i++ {
		r, eof := src.peekRune()
		if r != 0 || !eof {
			t.Fatalf("peekRune(); got:%d,%v, want:0,true", r, eof)
		}
	}}
