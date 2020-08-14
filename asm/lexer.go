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
type tokHash struct{}
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
type tokString struct {
	s string
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

type tokError struct {
	s       string
	source  *source
	lineNo  int
	linePos int
}

// lexer is an object that converts a stream of characters into a stream of tokens
type lexer struct {
	src *source
	nextToken token
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
func (l *lexer) getWord(firstRune rune) string {
	word := make([]rune, 1, 64)
	word[0] = firstRune
	for {
		r, eof := l.src.peekRune()
		if eof || (!unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_') {
			return string(word)
		}
		word = append(word, unicode.ToLower(r))
		l.src.consumeRune()
	}
}

// getAllowedString returns the next hex/bin/decimal digit string from the stream.
func (l *lexer) getAllowedString(firstRune rune, allowed string) string {
	str := make([]rune, 0, 32)
	if firstRune != 0 {
		str = append(str, firstRune)
	}
	for {
		r, eof := l.src.peekRune()
		r = unicode.ToLower(r)
		if eof || strings.IndexRune(allowed, r) == -1 {
			return string(str)
		}
		str = append(str, r)
		l.src.consumeRune()
	}
}

// getNumber returns a number from the stream. If firstRune > 0 it is the
// first rune of the number, which has already been consumed.
func (l *lexer) getIntNumber(firstRune rune, base int, allowed string) token {
	pos := l.src.curPos
	str := l.getAllowedString(firstRune, allowed)
	if str == "" {
		return &tokError{
			s:       "Illegal number (empty string)",
			source:  l.src,
			lineNo:  l.src.lineNo,
			linePos: pos,
		}
	}
	i, err := strconv.ParseInt(str, base, 64)
	if err != nil {
		return &tokError{
			s:       err.Error(),
			source:  l.src,
			lineNo:  l.src.lineNo,
			linePos: pos,
		}
	}
	return &tokIntNumber{i}
}

// getIdentifier takes a first rune, parses a valid identifier out of the stream
// and then returns the right token for it.
func (l *lexer) getIdentifier(firstRune rune) token {
	id := l.getWord(firstRune)
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

// getString returns a string parsed from the input stream.
func (l *lexer) getString() token {
	b := strings.Builder{}
	for {
		r, eof := l.src.peekRune()
		if r == '\n' || eof {
			return &tokError{
				s:       "Unexpected end-of-line",
				source:  l.src,
				lineNo:  l.src.lineNo,
				linePos: l.src.curPos,
			}
		}
		if r == '"' {
			l.src.consumeRune()
			r, _ := l.src.peekRune()
			if r != '"' {
				return &tokString{s: b.String()}
			}
		}
		b.WriteRune(r)
		l.src.consumeRune()
	}
}

// pushback pushes a single token back into the stream. If tok is nil
// it is not pushed back.
func (l *lexer) pushback(tok token) {
	if tok == nil {
		return
	}
	if l.nextToken != nil {
		panic(fmt.Sprintf("pushbackToken(%T %v) while nextToken=%F %v", tok, tok, l.nextToken, l.nextToken))
	}
	l.nextToken = tok
}

// getToken returns the next token in the stream.
func (l *lexer) getToken() token {
	if l.nextToken != nil {
		tok := l.nextToken
		l.nextToken = nil
		return tok
	}
	var r rune
	for {
		var eof bool
		r, eof = l.src.consumeRune()
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
		return &tokNewLine{}
	}
	if unicode.IsDigit(r) && r != '0' {
		// Special casing for 0x etc below.
		return l.getIntNumber(r, 10, decimalDigits)
	}
	r = unicode.ToLower(r)
	if unicode.IsLetter(r) || r == '_' {
		// Identifier or keyword.
		return l.getIdentifier(r)
	}
	switch r {
	case '0':
		// A number.
		r, eof := l.src.peekRune()
		if r == '\n' || eof {
			// A number just before the newline.
			return &tokIntNumber{0}
		}
		if r == 'x' {
			l.src.consumeRune()
			return l.getIntNumber(0, 16, hexDigits)
		}
		if r == 'b' {
			l.src.consumeRune()
			return l.getIntNumber(0, 2, binaryDigits)
		}
		return l.getIntNumber('0', 8, octalDigits)
	case '\'':
		r, _ := l.src.consumeRune()
		t := &tokRune{r}
		r, eof := l.src.peekRune()
		if !eof && r == '\'' {
			l.src.consumeRune()
			return t
		}
		return &tokError{
			s:       "Expected ' to end character constant",
			source:  l.src,
			lineNo:  l.src.lineNo,
			linePos: l.src.curPos,
		}
	case '"':
		return l.getString()
	case '$':
		return l.getIntNumber(0, 16, hexDigits)
	case '#':
		return &tokHash{}
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

