package jsonparser

import "errors"

func Parse(input string) (any, error) {
	l := lexer{}
	l.start(input)

	tok, ok := l.next()
	if !ok {
		return nil, errors.New("failed to read token")
	}

	return parseValue(tok, &l)
}
