package asm

import "fmt"

// context is the assembly context.
type context struct{
	pass int
	lexer *lexer
	seg *segment
	errors int
	warnings int
}

// lexError deals with the possibility of an error coming back from the lexer. Returns
// true if there was a lexer error (the error will already have been handled).
func (ctx *context) lexError(tok token) bool {
	if t, ok := tok.(*tokError); ok {
		fmt.Printf("[%d:%d] error: %s\n", t.lineNo, t.linePos, t.s)
		ctx.errors++
		return true
	}
	return false
}

// error reports an error.
func (ctx *context) error(s string, args ...interface{}) {
	fmt.Printf("[%d:%d] error: %s\n", ctx.lexer.src.lineNo, ctx.lexer.src.curPos, fmt.Sprintf(s, args...))
	ctx.errors++
}

// warning reports a warning.
func (ctx *context) warning(s string, args ...interface{}) {
	fmt.Printf("[%d:%d] warning: %s", ctx.lexer.src.lineNo, ctx.lexer.src.curPos, fmt.Sprintf(s, args...))
	ctx.warnings++
}

// expect expects a token and registers an error if that token did not appear,
func (ctx *context) expect(f func(token) bool, typ string) (tok token, ok bool) {
	tok = ctx.lexer.getToken()
	ok = f(tok)
	if !ok {
		ctx.error("expected %s, not '%T'", typ, tok)
	}
	return tok, ok
}
