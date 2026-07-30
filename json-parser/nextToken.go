package jsonparser

import "errors"

func nextToken(l *lexer) (token, error) {
	tok, ok := l.next()

	if !ok {
		return token{}, errors.New("unexpected end of file")
	}

	return tok, nil
}
