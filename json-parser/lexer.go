package jsonparser

type lexer struct {
	position int
	input    string
}

func (l *lexer) start(input string) {
	l.input = input
	l.position = 0
}

func (l *lexer) readString() (string, bool) {
	ok := false
	startingPosition := l.position
	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		// add skipping backslash

		if char == '"' {
			ok = true
			break
		}
	}

	if !ok {
		return "", false
	}

	return l.input[startingPosition : l.position-1], true
}

// next returns the next token. ok is false at EOF.
func (l *lexer) next() (token, bool) {
	// var inString bool

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

	}
	return token{
		kind: tokenEOF,
		lit:  "",
	}, false
}
