package jsonparser

import "errors"

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
			stringInternal, err := l.readString()
			return tok(tokenString, stringInternal), err
		case ' ', '\t', '\n', '\r':
			continue
		case ':':
			return tok(tokenColon, ":"), nil
		case ',':
			return tok(tokenComma, ","), nil
		case '{':
			return tok(tokenLBrace, "{"), nil
		case '}':
			return tok(tokenRBrace, "}"), nil
		case '[':
			return tok(tokenLBracket, "["), nil
		case ']':
			return tok(tokenRBracket, "]"), nil
		}

		if (char >= '0' && char <= '9') || char == '-' {
			result, err := l.readNumber()

			if err != nil {
				return token{}, err
			}

			return tok(tokenNumber, result), nil
		}

		if char >= 'a' && char <= 'z' {
			return l.readKeyword()
		}

		return token{}, errors.New("unrecognized token")

	}
	return tok(tokenEOF, ""), nil
}
