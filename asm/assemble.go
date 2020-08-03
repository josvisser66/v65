// Package asm contains all the code required to assemble a source file.
package asm

import "fmt"
import "errors"

// lineStarter is an interface that tokens implement if they can start a line.
type lineStarter interface {
	assemble(seg *segment, src *source, tok token) *source
}

// Assemble assembles a source file.
func Assemble(filename string) (*segment, error) {
	src, err := newSource(filename)
	if err != nil {
		return nil, err
	}
	return assemble(src)
}

// assemble assembles from a source object.
func assemble(src *source) (*segment, error) {
	segment := newSegment()

loop:
	for {
		if tok := src.getToken(); segment.lexError(tok) {
			src.skipRestOfLine()
		} else {
			switch tok.(type) {
			case lineStarter:
				src = tok.(lineStarter).assemble(segment, src, tok)
			case *tokEOF:
				break loop
			case *tokNewLine:
			default:
				segment.error(src, "unexpected token at start of line: %T", tok)
				src.skipRestOfLine()
			}
		}
	}

	fmt.Sprintf("There were %d error(s) and %d warning(s).\n", segment.errors, segment.warnings)

	if segment.errors > 0 {
		return nil, errors.New("no output generated because of errors")
	}

	return segment, nil
}
