package jsonparser

import (
	"fmt"
	"strconv"
)

func parseValue(tok token, l *lexer) (any, error) {
	switch tok.kind {
	case tokenLBrace:
		return parseObject(l)
	case tokenLBracket:
		return parseArray(l)
	case tokenString:
		return tok.lit, nil
	case tokenNumber:
		num, err := strconv.ParseFloat(tok.lit, 64)

		if err != nil {
			return nil, fmt.Errorf("invalid number %q at %d: %w", tok.lit, tok.pos, err)
		}

		return num, nil
	case tokenNull:
		return nil, nil
	case tokenTrue:
		return true, nil
	case tokenFalse:
		return false, nil
	}

	return nil, fmt.Errorf("expected value, found %q at %d", tok.lit, tok.pos)
}
