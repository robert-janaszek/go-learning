package jsonparser

import (
	"errors"
	"strings"
)

func (l *lexer) readString() (string, error) {
	var result strings.Builder
	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		switch char {
		case '"':
			return result.String(), nil
		case '\\':
			if l.position >= len(l.input) {
				return "", errors.New("unexpected EOF, unterminated escape")
			}
			result.WriteByte(l.input[l.position])
			l.position++
		default:
			result.WriteByte(char)
		}

	}

	return "", errors.New(`unexpected EOF, wanted '"'`)
}
