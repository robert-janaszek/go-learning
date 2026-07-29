package jsonparser

import "errors"

func parseArray(l *lexer) (any, error) {
	var values []any = []any{}

	var tok token
	var ok bool

	tok, ok = l.next()

	if !ok {
		return nil, errors.New("unexpected end of file")
	}

	if tok.kind == tokenRBracket {
		return values, nil
	}

	var val any
	var err error
	val, err = parseValue(tok, l)

	if err != nil {
		return nil, err
	}

	values = append(values, val)

	tok, ok = l.next()
	if !ok {
		return nil, errors.New("unexpected end of file")
	}

	for {
		if tok.kind == tokenRBracket {
			break
		}

		if tok.kind != tokenComma {
			return nil, errors.New("expected comma")
		}

		tok, ok = l.next()
		if !ok {
			return nil, errors.New("unexpected end of file")
		}

		val, err := parseValue(tok, l)

		if err != nil {
			return nil, err
		}

		values = append(values, val)

		tok, ok = l.next()

		if !ok {
			return nil, errors.New("unexpected end of file")
		}
	}

	return values, nil
}
