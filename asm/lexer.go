package asm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type token int

const (
	hexDigits = "0123456789abcdef"
	binaryDigits = "01"
	octalDigits = "01234567"
	decimalDigits = "01234567890"

	tokEOF token = iota
	tokNewLine
	tokIdentifier
	tokNumber
	tokChar
	tokError
)

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
func (s *source) getNumber(firstRune rune, base int, allowed string) (t token, text string, val int64, eof bool) {
	str := s.getString(firstRune, allowed)
	if str == "" {
		return tokError, "illegal number <empty string>", 0, false
	}
	i, err := strconv.ParseInt(str, base, 64)
	if err != nil {
		return tokError, fmt.Sprintf("number too large (base %d): %s", base, str), 0, false
	}
	return tokNumber, fmt.Sprintf("%s", str), i, false
}

// getToken returns the next token in the stream.
func (s *source) getToken() (t token, text string, val int64, eof bool) {
	var r rune
	for {
		var eof bool
		r, eof = s.consumeRune()
		if eof {
			return tokEOF, "", 0, eof
		}
		if r == '\n' {
			return tokNewLine, "\n", 0, false
		}
		if !unicode.IsSpace(r) {
			break
		}
	}
	if r == ';' {
		// Ignores comments.
		s.skipToEOLN()
		return tokNewLine, "\n", 0, false
	}
	r = unicode.ToLower(r)
	if unicode.IsLetter(r) || r == '_' {
		// Identifier or keyword.
		return tokIdentifier, s.getWord(r), 0, false
	}

	if r == '0' {
		// A number.
		r, eof = s.peekRune()
		if r == '\n' || eof {
			// A number just before the newline.
			return tokNumber, "0", 0, false
		}
		if r == 'x' {
			s.consumeRune()
			return s.getNumber(0, 16, hexDigits)
		}
		if r == 'b' {
			s.consumeRune()
			return s.getNumber(0, 2, binaryDigits)
		}
		return s.getNumber('0', 8, octalDigits)
	}
	if unicode.IsDigit(r) {
		return s.getNumber(r, 10, decimalDigits)
	}
	if r == '$' {
		return s.getNumber(0, 16, hexDigits)
	}
	return tokChar, string(r), 0, false
}
