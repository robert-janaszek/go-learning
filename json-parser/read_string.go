package jsonparser

import (
	"fmt"
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

				return "", fmt.Errorf("unexpected EOF, unterminated escape at %d", l.position-1)
			}
			result.WriteByte(l.input[l.position])
			l.position++
		default:
			result.WriteByte(char)
		}

	}

	return "", fmt.Errorf(`unexpected EOF, wanted '"' at %d`, l.position-1)
}
