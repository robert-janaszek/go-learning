package jsonparser

import "strings"

func (l *lexer) readString() (string, bool) {
	var result strings.Builder
	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		switch char {
		case '"':
			return result.String(), true
		case '\\':
			if l.position >= len(l.input) {
				return "", false
			}
			result.WriteByte(l.input[l.position])
			l.position++
		default:
			result.WriteByte(char)
		}

	}

	return "", false
}
