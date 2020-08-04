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
	// with an external symbol and then a positive or negative constant
	// offset.
	if id, ok := tok.(*tokIdentifier); ok {
		var sym symbol
		if sym, ok = seg.symbols[id.id]; !ok {
			// The error will be generated down there somewhere.
			v, next := s.level1(seg, tok)
			return &exprValue{nil, v}, next
		}
		if externSym, ok := sym.(*externSymbol); ok {
			// This is an external symbol. The rest of the expression
			// can be + or - something else.
			next := s.getToken()
			var v int64
			if _, ok := next.(*tokPlus); ok {
				v, next = s.level1(seg, nil)
			} else if _, ok := next.(*tokMinus); ok {
				v, next = s.level1(seg, nil)
			}
			return &exprValue{externSym, v}, next
		}
		// The identifier is a label with a known value. Fallthrough.
	}

	v, next := s.level1(seg, tok)
	return &exprValue{nil, v}, next
}

func (s *source) level1(seg *segment, nextToken token) (val int64, next token) {
	val, next = s.level2(seg, nextToken)
	var v int64
	for {
		if next == nil {
			next = s.getToken()
		}
		if _, ok := next.(*tokOr); ok {
			v, next = s.level2(seg, nil)
			val = val | v
		} else if _, ok := next.(*tokAnd); ok {
			v, next = s.level2(seg, nil)
			val = val & v
		} else {
			return val, next
		}
	}
}

func (s *source) level2(seg *segment, nextToken token) (val int64, next token) {
	val, next = s.level3(seg, nextToken)
	var v int64
	for {
		if next == nil {
			next = s.getToken()
		}
		if _, ok := next.(*tokPlus); ok {
			v, next = s.level3(seg, nil)
			val = val + v
		} else if _, ok := next.(*tokMinus); ok {
			v, next = s.level3(seg, nil)
				val = val - v
		} else {
			return val, next
		}
	}
}

func (s *source) level3(seg *segment, nextToken token) (val int64, next token) {
	val, next = s.level4(seg, nextToken)
	var v int64
	for {
		if next == nil {
			next = s.getToken()
		}
		if _, ok := next.(*tokMultiply); ok {
			v, next = s.level4(seg, nil)
			val = val * v
		} else if _, ok := next.(*tokDivide); ok {
			v, next = s.level4(seg, nil)
			if v == 0 {
				seg.error(s, "division by zero")
			} else {
				val = val / v
			}
		} else {
			return val, next
		}
	}
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
		v, next := s.level1(seg, nil)
		if _, ok := next.(*tokRightParen); !ok {
			seg.error(s, "expected ')', not '%T'", next)
			s.skipToEOLN()
			return 0, nil
		}
		return v , nil
	}
	if id, ok := nextToken.(*tokIdentifier); ok {
		// Label, must be locally defined.
		sym, ok := seg.symbols[id.id]
		if !ok {
			seg.error(s, "unknown label: %s", id.id)
			return 0, nil
		}
		if localSym, ok := sym.(*localSymbol); ok {
			return localSym.value, nil
		}
		seg.error(s, "illegal use in expression of external label: %s", id.id)
		return 0, nil
	}
	if _, ok := nextToken.(*tokPlus); ok {
		// Unary plus operator.
		v, next := s.level4(seg, nil)
		return v, next
	}
	if _, ok := nextToken.(*tokMinus); ok {
		// Unary minus operator.
		v, next := s.level4(seg, nil)
		return -v, next
	}
	seg.error(s, "invalid expression; unexpected token: '%T'", nextToken)
	s.skipRestOfLine()
	return 0, nil
}