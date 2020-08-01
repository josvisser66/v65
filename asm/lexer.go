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
type tokError struct {
	s       string
	source  *source
	lineNo  int
	linePos int
}

// Error returns a string description of this error.
// This function makes tokError an error.
func (t *tokError) Error() string {
	return fmt.Sprintf("%d:%d error: %s", t.lineNo, t.linePos, t.s)
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
func (s *source) getIntNumber(firstRune rune, base int, allowed string) (token, error) {
	pos := s.curPos
	str := s.getString(firstRune, allowed)
	if str == "" {
		return nil, &tokError{
			s:       "Illegal number (empty string)",
			source:  s,
			lineNo:  s.lineNo,
			linePos: pos,
		}
	}
	i, err := strconv.ParseInt(str, base, 64)
	if err != nil {
		return nil, &tokError{
			s:       err.Error(),
			source:  s,
			lineNo:  s.lineNo,
			linePos: pos,
		}
	}
	return &tokIntNumber{i}, nil
}

// getToken returns the next token in the stream.
func (s *source) getToken() (token, error) {
	var r rune
	for {
		var eof bool
		r, eof = s.consumeRune()
		if eof {
			return &tokEOF{}, nil
		}
		if r == '\n' {
			return &tokNewLine{}, nil
		}
		if !unicode.IsSpace(r) {
			break
		}
	}
	if r == ';' {
		// Ignores comments.
		s.skipToEOLN()
		return &tokNewLine{}, nil
	}
	r = unicode.ToLower(r)
	if unicode.IsLetter(r) || r == '_' {
		// Identifier or keyword.
		return &tokIdentifier{s.getWord(r)}, nil
	}

	if r == '0' {
		// A number.
		r, eof := s.peekRune()
		if r == '\n' || eof {
			// A number just before the newline.
			return &tokIntNumber{0}, nil
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
	}
	if unicode.IsDigit(r) {
		return s.getIntNumber(r, 10, decimalDigits)
	}
	if r == '$' {
		return s.getIntNumber(0, 16, hexDigits)
	}
	return &tokRune{r}, nil
}
