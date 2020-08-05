// Package asm contains all the code required to assemble a source file.
package asm

import "fmt"

// lineStarter is an interface that tokens implement if they can start a line.
type lineStarter interface {
	assemble(ctx *context)
}

// Assemble assembles a source file.
func Assemble(filename string) (*context, error) {
	src, err := newSource(filename)
	if err != nil {
		return nil, err
	}
	ctx := &context{pass: 1, seg: newSegment(), src: src}
	ctx.assemble()
	if ctx.errors == 0 {
		ctx.pass++
		ctx.seg.lc = 0
		ctx.assemble()
	}
	fmt.Sprintf("There were %d error(s) and %d warning(s).\n", ctx.errors, ctx.warnings)
	return ctx, nil
}

// assemble assembles from a source object.
func (ctx *context) assemble() {
loop:
	for {
		if tok := ctx.src.getToken(); ctx.lexError(tok) {
			ctx.src.skipRestOfLine()
		} else {
			switch tok.(type) {
			case lineStarter:
				tok.(lineStarter).assemble(ctx)
			case *tokEOF:
				break loop
			case *tokNewLine:
			default:
				ctx.error("unexpected token at start of line: %T", tok)
				ctx.src.skipRestOfLine()
			}
		}
	}
}
