package jsonparser

import "errors"

func (l *lexer) readKeyword() (token, error) {
	startingPosition := l.position - 1

	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		if char >= 'a' && char <= 'z' {
			continue
		}

		l.position--
		break
	}

	possibleKeyword := l.input[startingPosition:l.position]

	switch possibleKeyword {
	case "null":
		return tok(tokenNull, "null", startingPosition), nil
	case "true":
		return tok(tokenTrue, "true", startingPosition), nil
	case "false":
		return tok(tokenFalse, "false", startingPosition), nil
	}

	return token{}, errors.New("unrecognized keyword: " + possibleKeyword)
}
