package jsonparser

import "strings"

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

func (l *lexer) readNumber(startChar byte) (string, bool) {
	eConsumed := false
	includeLast := true
	startingPosition := l.position

	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		if char >= '0' && char <= '9' {
			continue
		}

		switch char {
		case 'e', 'E':
			if eConsumed {
				return "", false
			}

			eConsumed = true
			continue
		}

		includeLast = false
		break
	}

	var stringRest string
	if includeLast {
		stringRest = l.input[startingPosition:l.position]
	} else {
		stringRest = l.input[startingPosition : l.position-1]
	}

	var resultBuilder strings.Builder
	resultBuilder.WriteByte(startChar)
	resultBuilder.WriteString(stringRest)

	return resultBuilder.String(), true
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

		if char >= '0' && char <= '9' {
			result, ok := l.readNumber(char)

			if !ok {
				return token{}, false
			}

			return token{
				kind: tokenNumber,
				lit:  result,
			}, true
		}

	}
	return token{
		kind: tokenEOF,
		lit:  "",
	}, false
}
