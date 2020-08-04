package asm

type exprValue struct {
	sym *externSymbol // If this is the value of a relocatable expression.
	val int64
}

func (s *source) expr(seg *segment) (val *exprValue, next token) {
	// Gets the token that starts the expression.
	tok := s.getToken()

	// If this token is an identifier, there is a change that we have
	// an expression that requires relocation. Such an expression starts
	// with an external symbol and then a positive or negative offset.
	if id, ok := tok.(*tokIdentifier); ok {
		var sym symbol
		if sym, ok = seg.symbols[id.id]; !ok {
			seg.error(s, "unknown label: %s", id.id)
			v, next := s.level1(seg, s.getToken())
			return &exprValue{nil, v}, next
		}
		if externSym, ok := sym.(*externSymbol); ok {
			// This is an external symbol.
			v, next := s.level1(seg, s.getToken())
			return &exprValue{externSym, v}, next
		}
		// The identifier is a label with a known value. Fallthrough.
	}

	v, next := s.level1(seg, tok)
	return &exprValue{nil, v}, next
}

func (s *source) level1(seg *segment, nextToken token) (val int64, next token) {
	val, next = s.level2(seg, nextToken)

	return val, next
}

func (s *source) level2(seg *segment, nextToken token) (val int64, next token) {
	val, next = s.level3(seg, nextToken)

	return val, next
}

func (s *source) level3(seg *segment, nextToken token) (val int64, next token) {
	val, next = s.level4(seg, nextToken)

	return val, next
}

func (s *source) level4(seg *segment, nextToken token) (val int64, next token) {
	if nextToken == nil {
		nextToken = s.getToken()
	}
	if seg.lexError(nextToken) {
		s.skipRestOfLine()
		return 0,nil
	}
	if num, ok := nextToken.(*tokIntNumber); ok {
		return num.n, s.getToken()
	}
	if _, ok := nextToken.(*tokLeftParen); ok {
		val, next := s.level1(seg, nil)
		if _, ok := next.(*tokRightParen); !ok {
			seg.error(s, "expected ')', not '%T'", next)
			s.skipToEOLN()
			return 0, nil
		}
		return val , nil
	}
	seg.error(s, "invalid expression; unexpected token: '%T'", nextToken)
	s.skipRestOfLine()
	return 0, nil
}