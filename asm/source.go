package asm

import (
	"io/ioutil"
	"os"
	"strings"
)

type source struct {
	lines    []string
	lineNo   int
	curLine  []rune
	curPos   int
	nextChar rune
}

// newSourceFromString creates a new source file struct from a string value.
// The data string will be split on newlines.
func newSourceFromString(data string) *source {
	return &source{lines: strings.Split(data, "\n"), lineNo: -1}
}

// newSource creates a new source file struct a file name.
func newSource(filename string) (*source, error) {
	var content []byte
	var err error
	if filename == "-" {
		content, err = ioutil.ReadAll(os.Stdin)
	} else {
		content, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		return nil, err
	}
	return newSourceFromString(string(content)), nil
}

// peekRune returns the next character without consuming it.
func (s *source) peekRune() (r rune, eof bool) {
	if s.nextChar != 0 {
		return s.nextChar, false
	}
	if s.curLine == nil || s.curPos > len(s.curLine) {
		if s.lineNo == len(s.lines) - 1 {
			return 0, true
		}
		s.lineNo++
		s.curLine = []rune(s.lines[s.lineNo])
		s.curPos = 0
	}
	if s.curPos == len(s.curLine) {
		s.nextChar = '\n'
	} else {
		s.nextChar = s.curLine[s.curPos]
	}
	s.curPos++
	return s.nextChar, false
}

// consumeRune consumes a character from the source (and returns it).
func (s *source) consumeRune() (r rune, eof bool) {
	r, eof = s.peekRune()
	s.nextChar = 0
	return
}
