package jsonparser

type lexer struct {
	position int
	input    string
}

func (l *lexer) start(input string) {
	l.input = input
	l.position = 0
}

// next returns the next token. ok is false at EOF.
func (l *lexer) next() (token, bool) {
	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		switch char {
		case '"':
			stringInternal, ok := l.readString()
			return tok(tokenString, stringInternal), ok
		case ' ', '\t', '\n', '\r':
			continue
		case ':':
			return tok(tokenColon, ":"), true
		case ',':
			return tok(tokenComma, ","), true
		case '{':
			return tok(tokenLBrace, "{"), true
		case '}':
			return tok(tokenRBrace, "}"), true
		case '[':
			return tok(tokenLBracket, "["), true
		case ']':
			return tok(tokenRBracket, "]"), true
		}

		if (char >= '0' && char <= '9') || char == '-' {
			result, ok := l.readNumber()

			if !ok {
				return token{}, false
			}

			return tok(tokenNumber, result), true
		}

		if char >= 'a' && char <= 'z' {
			return l.readKeyword()
		}

		return token{}, false

	}
	return tok(tokenEOF, ""), false
}
