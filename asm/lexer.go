package asm

import "unicode"

type token int

const (
	tokEOF token = iota
	tokNewLine
	tokIdentifier
	tokUnknown
)

// getWord returns the next word from the stream. firstRune is the first
// (starter) rune for the word, which has already been consumed from
// the string.
func (s * source) getWord(firstRune rune) string {
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

// getToken returns the next token in the stream.
func (s *source) getToken() (t token, text string, eof bool) {
	var r rune
	for {
		var eof bool
		r, eof = s.consumeRune()
		if eof {
			return tokEOF, "", eof
		}
		if r == '\n' {
			return tokNewLine, "\n", false
		}
		if !unicode.IsSpace(r) {
			break
		}
	}
	if r == ';' {
		// Ignores comments.
		s.skipToEOLN()
		return tokNewLine, "\n", false
	}
	r = unicode.ToLower(r)
	if unicode.IsLetter(r) || r == '_' {
		// Identifier or keyword.
		return tokIdentifier, s.getWord(r), false
	}
	return tokUnknown, string(r), false
}
