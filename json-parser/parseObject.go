package jsonparser

import (
	"fmt"
)

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

	err = expectToken(tok, tokenString)
	if err != nil {
		return nil, err
	}
	foundKey = tok.lit

	tok, err = nextToken(l)
	if err != nil {
		return nil, err
	}

	err = expectToken(tok, tokenColon)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("expected ',' or '}', found %q at %d", tok.lit, tok.pos)
		}

		tok, err = nextToken(l)

		if err != nil {
			return nil, err
		}

		err = expectToken(tok, tokenString)
		if err != nil {
			return nil, err
		}

		foundKey = tok.lit

		tok, err = nextToken(l)

		if err != nil {
			return nil, err
		}

		err = expectToken(tok, tokenColon)
		if err != nil {
			return nil, err
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
