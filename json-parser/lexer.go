package jsonparser

import (
	"fmt"
)

type lexer struct {
	position int
	input    string
}

func (l *lexer) start(input string) {
	l.input = input
	l.position = 0
}

// next returns the next token. ok is false at EOF.
func (l *lexer) next() (token, error) {
	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		switch char {
		case '"':
			pos := l.position - 1
			stringInternal, err := l.readString(pos)
			return tok(tokenString, stringInternal, pos), err
		case ' ', '\t', '\n', '\r':
			continue
		case ':':
			return tok(tokenColon, ":", l.position-1), nil
		case ',':
			return tok(tokenComma, ",", l.position-1), nil
		case '{':
			return tok(tokenLBrace, "{", l.position-1), nil
		case '}':
			return tok(tokenRBrace, "}", l.position-1), nil
		case '[':
			return tok(tokenLBracket, "[", l.position-1), nil
		case ']':
			return tok(tokenRBracket, "]", l.position-1), nil
		}

		if (char >= '0' && char <= '9') || char == '-' {
			pos := l.position - 1
			result, err := l.readNumber(char)

			if err != nil {
				return token{}, err
			}

			return tok(tokenNumber, result, pos), nil
		}

		if char >= 'a' && char <= 'z' {
			return l.readKeyword()
		}

		return token{}, fmt.Errorf("unrecognized character %q at %d", char, l.position-1)

	}
	return tok(tokenEOF, "", l.position), nil
}
