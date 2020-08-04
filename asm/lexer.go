package asm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	hexDigits     = "0123456789abcdef"
	binaryDigits  = "01"
	octalDigits   = "01234567"
	decimalDigits = "01234567890"
)

type token interface {
}
type tokEOF struct{}
type tokComma struct{}
type tokLeftParen struct{}
type tokRightParen struct{}
type tokOr struct{}
type tokAnd struct{}
type tokPlus struct{}
type tokMinus struct{}
type tokMultiply struct{}
type tokDivide struct{}
type tokNewLine struct{}
type tokIdentifier struct {
	id string
}
type tokIntNumber struct {
	n int64
}
type tokRune struct {
	r rune
}
type tokOpcode struct {
	opcode string
}
type tokRegisterA struct{}
type tokRegisterX struct{}
type tokRegisterY struct{}

type tokEqu struct{}
type tokInclude struct{}
type tokError struct {
	s       string
	source  *source
	lineNo  int
	linePos int
}

// metaMap is a map from an identifier to a token that represents
// an assembler meta instruction.
var metaMap = make(map[string]token)

// Error returns a string description of this error.
// This function makes tokError an error.
func (t *tokError) Error() string {
	return fmt.Sprintf("[%d:%d] %s", t.lineNo, t.linePos, t.s)
}

// getWord returns the next word from the stream. firstRune is the first
// (starter) rune for the word, which has already been consumed from
// the string.
func (s *source) getWord(firstRune rune) string {
	word := make([]rune, 1, 64)
	word[0] = firstRune
	for {
		r, eof := s.peekRune()
		if eof || (!unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_') {
			return string(word)
		}
		word = append(word, unicode.ToLower(r))
		s.consumeRune()
	}
}

// getHexString returns the next hex/bin/decimal digit string from the stream.
func (s *source) getString(firstRune rune, allowed string) string {
	str := make([]rune, 0, 32)
	if firstRune != 0 {
		str = append(str, firstRune)
	}
	for {
		r, eof := s.peekRune()
		r = unicode.ToLower(r)
		if eof || strings.IndexRune(allowed, r) == -1 {
			return string(str)
		}
		str = append(str, r)
		s.consumeRune()
	}
}

// getNumber returns a number from the stream. If firstRune > 0 it is the
// first rune of the number, which has already been consumed.
func (s *source) getIntNumber(firstRune rune, base int, allowed string) token {
	pos := s.curPos
	str := s.getString(firstRune, allowed)
	if str == "" {
		return &tokError{
			s:       "Illegal number (empty string)",
			source:  s,
			lineNo:  s.lineNo,
			linePos: pos,
		}
	}
	i, err := strconv.ParseInt(str, base, 64)
	if err != nil {
		return &tokError{
			s:       err.Error(),
			source:  s,
			lineNo:  s.lineNo,
			linePos: pos,
		}
	}
	return &tokIntNumber{i}
}

// getIdentifier takes a first rune, parses a valid identifier out of the stream
// and then returns the right token for it.
func (s *source) getIdentifier(firstRune rune) token {
	id := s.getWord(firstRune)
	if _, ok := opcodes[id]; ok {
		return &tokOpcode{opcode: id}
	}
	if tok, ok := metaMap[id]; ok {
		return tok
	}
	switch {
	case id == "a":
		return &tokRegisterA{}
	case id == "x":
		return &tokRegisterX{}
	case id == "y":
		return &tokRegisterY{}
	}
	return &tokIdentifier{id: id}
}

// getToken returns the next token in the stream.
func (s *source) getToken() token {
	var r rune
	for {
		var eof bool
		r, eof = s.consumeRune()
		if eof {
			return &tokEOF{}
		}
		if r == '\n' {
			return &tokNewLine{}
		}
		if !unicode.IsSpace(r) {
			break
		}
	}
	if r == ';' {
		// Ignores comments.
		s.skipRestOfLine()
		return &tokNewLine{}
	}
	if unicode.IsDigit(r) && r != '0' {
		// Special casing for 0x etc below.
		return s.getIntNumber(r, 10, decimalDigits)
	}
	r = unicode.ToLower(r)
	if unicode.IsLetter(r) || r == '_' {
		// Identifier or keyword.
		return s.getIdentifier(r)
	}
	switch r {
	case '0':
		// A number.
		r, eof := s.peekRune()
		if r == '\n' || eof {
			// A number just before the newline.
			return &tokIntNumber{0}
		}
		if r == 'x' {
			s.consumeRune()
			return s.getIntNumber(0, 16, hexDigits)
		}
		if r == 'b' {
			s.consumeRune()
			return s.getIntNumber(0, 2, binaryDigits)
		}
		return s.getIntNumber('0', 8, octalDigits)
	case '$':
		return s.getIntNumber(0, 16, hexDigits)
	case ',':
		return &tokComma{}
	case '|':
		return &tokOr{}
	case '&':
		return &tokAnd{}
	case '+':
		return &tokPlus{}
	case '-':
		return &tokMinus{}
	case '*':
		return &tokMultiply{}
	case '/':
		return &tokDivide{}
	case '(':
		return &tokLeftParen{}
	case ')':
		return &tokRightParen{}
	default:
		return &tokRune{r}
	}
}

func (s *source) expect(seg *segment, f func(token) bool, typ string) (tok token, ok bool) {
	tok = s.getToken()
	ok = f(tok)
	if !ok {
		seg.error(s, "expected %s, not '%T'", typ, tok)
	}
	return tok, ok
}
