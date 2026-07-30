package jsonparser

import (
	"fmt"
)

func Parse(input string) (any, error) {
	l := lexer{}
	l.start(input)

	tok, err := l.next()
	if err != nil {
		return nil, err
	}

	val, err := parseValue(tok, &l)

	if err != nil {
		return nil, err
	}

	tok, err = l.next()

	if err != nil {
		return nil, err
	}

	if tok.kind != tokenEOF {
		return nil, fmt.Errorf("wanted EOF, found: %q", tok.lit)
	}

	return val, nil
}
