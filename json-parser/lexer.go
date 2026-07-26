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
			return token{
				kind: tokenString,
				lit:  stringInternal,
			}, ok
		case ' ', '\t', '\n', '\r':
			continue
		case ':':
			return token{
				kind: tokenColon,
				lit:  ":",
			}, true
		case ',':
			return token{
				kind: tokenComma,
				lit:  ",",
			}, true
		case '{':
			return token{
				kind: tokenLBrace,
				lit:  "{",
			}, true
		case '}':
			return token{
				kind: tokenRBrace,
				lit:  "}",
			}, true
		case '[':
			return token{
				kind: tokenLBracket,
				lit:  "[",
			}, true
		case ']':
			return token{
				kind: tokenRBracket,
				lit:  "]",
			}, true
		}

		if (char >= '0' && char <= '9') || char == '-' {
			result, ok := l.readNumber()

			if !ok {
				return token{}, false
			}

			return token{
				kind: tokenNumber,
				lit:  result,
			}, true
		}

		if char >= 'a' && char <= 'z' {
			return l.readKeyword()
		}

		return token{}, false

	}
	return token{
		kind: tokenEOF,
		lit:  "",
	}, false
}
