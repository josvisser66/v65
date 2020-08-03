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
		tok := src.getToken()
		if segment.lexError(tok) {
			src.skipToEOLN()
			continue
		}
		if ls, ok := tok.(lineStarter); ok {
			src = ls.assemble(segment, src, tok)
			continue

		}
		switch tok.(type) {
		case *tokEOF:
			break loop
		case *tokNewLine:
		default:
			segment.error(src, "unexpected token: %T", tok)
			src.skipToEOLN()
		}
	}

	fmt.Sprintf("There were %d error(s) and %d warning(s).\n", segment.errors, segment.warnings)

	if segment.errors > 0 {
		return nil, errors.New("no output generated because of errors")
	}

	return segment, nil
}
