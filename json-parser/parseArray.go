package jsonparser

func parseArray(l *lexer) (any, error) {
	var values []any = []any{}

	tok, err := nextToken(l)
	if err != nil {
		return nil, err
	}

	if tok.kind == tokenRBracket {
		return values, nil
	}

	var val any
	val, err = parseValue(tok, l)

	if err != nil {
		return nil, err
	}

	values = append(values, val)

	tok, err = nextToken(l)
	if err != nil {
		return nil, err
	}

	for {
		if tok.kind == tokenRBracket {
			break
		}

		err = expectToken(tok, tokenComma)
		if err != nil {
			return nil, err
		}

		tok, err = nextToken(l)
		if err != nil {
			return nil, err
		}

		val, err := parseValue(tok, l)

		if err != nil {
			return nil, err
		}

		values = append(values, val)

		tok, err = nextToken(l)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}
