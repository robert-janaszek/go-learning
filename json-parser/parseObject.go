package jsonparser

import "errors"

func parseObject(l *lexer) (any, error) {
	obj := map[string]any{}
	var tok token
	var err error
	var val any
	var foundKey string

	tok, err = nextToken(l)
	if err != nil {
		return nil, err
	}

	if tok.kind == tokenRBrace {
		return obj, nil
	}

	if tok.kind != tokenString {
		return nil, errors.New("expected string, but found " + tok.lit)
	}
	foundKey = tok.lit

	tok, err = nextToken(l)
	if err != nil {
		return nil, err
	}

	if tok.kind != tokenColon {
		return nil, errors.New("expected ':', but found " + tok.lit)
	}

	tok, err = nextToken(l)
	if err != nil {
		return nil, err
	}

	val, err = parseValue(tok, l)

	if err != nil {
		return nil, err
	}

	obj[foundKey] = val

	for {
		tok, err = nextToken(l)
		if err != nil {
			return nil, err
		}

		if tok.kind == tokenRBrace {
			return obj, nil
		}

		if tok.kind != tokenComma {
			return nil, errors.New("expected ',' or '}', but found: " + tok.lit)
		}

		tok, err = nextToken(l)

		if err != nil {
			return nil, err
		}

		if tok.kind != tokenString {
			return nil, errors.New("expected string, but found " + tok.lit)
		}

		foundKey = tok.lit

		tok, err = nextToken(l)

		if err != nil {
			return nil, err
		}

		if tok.kind != tokenColon {
			return nil, errors.New("expected ':', but found " + tok.lit)
		}

		tok, err = nextToken(l)

		if err != nil {
			return nil, err
		}

		val, err = parseValue(tok, l)

		if err != nil {
			return nil, err
		}

		obj[foundKey] = val
	}
}
