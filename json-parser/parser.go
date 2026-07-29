package jsonparser

import "errors"

func Parse(input string) (any, error) {
	l := lexer{}
	l.start(input)

	tok, ok := l.next()
	if !ok {
		return nil, errors.New("failed to read token")
	}

	val, err := parseValue(tok, &l)

	if err != nil {
		return nil, err
	}

	tok, ok = l.next()

	if ok {
		return nil, errors.New("expected EOF, but found " + tok.lit)
	}

	return val, nil
}
