package asm

import (
	"io/ioutil"
	"os"
	"strings"
)

type source struct {
	filename string
	lines    []string
	lineNo   int // 1-based, so not directly an index into lines
	curLine  []rune
	curPos   int // 1-based, so not directly an index into curLine
	nextChar rune
}

// newSourceFromString creates a new source file struct from a string value.
// The data string will be split on newlines.
func newSourceFromString(data string) *source {
	return &source{lines: strings.Split(data, "\n"), lineNo: 0}
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
	s := newSourceFromString(string(content))
	s.filename = filename
	return s, nil
}

// peekRune returns the next character without consuming it.
func (s *source) peekRune() (r rune, eof bool) {
	if s.nextChar != 0 {
		return s.nextChar, false
	}
	// If this is the first time we are calling peekRune, we move
	// to the next line.
	if s.lineNo == 0 {
		s.moveToNextLine()
	}
	// If we are pointed at the line after the last one, we return eof.
	if s.lineNo == len(s.lines)+1 {
		return 0, true
	}
	// If we exhausted this line we keep returning newlines.
	if s.curPos == len(s.curLine) + 1 {
		s.nextChar = '\n'
	} else {
		s.nextChar = s.curLine[s.curPos-1]
		s.curPos++
	}
	return s.nextChar, false
}

// consumeRune consumes a character from the source (and returns it).
func (s *source) consumeRune() (r rune, eof bool) {
	r, eof = s.peekRune()
	s.nextChar = 0
	return
}

// moveToNextLine moves to the next line of the input. We never move
// beyond the line after the next line though.
func (s *source) moveToNextLine() {
	s.nextChar = 0
	if s.lineNo == len(s.lines) + 1 {
		return
	}
	s.lineNo++
	s.curLine = nil
	s.curPos = 1
	if s.lineNo <= len(s.lines) {
		s.curLine = []rune(s.lines[s.lineNo-1])
	}
}
