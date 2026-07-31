package jsonparser

import (
	"fmt"
	"strings"
)

func (l *lexer) readString(startPos int) (string, error) {
	var result strings.Builder
	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		switch char {
		case '"':
			return result.String(), nil
		case '\\':
			if l.position >= len(l.input) {
				return "", fmt.Errorf("unterminated escape at %d", startPos)
			}
			result.WriteByte(l.input[l.position])
			l.position++
		default:
			result.WriteByte(char)
		}
	}

	return "", fmt.Errorf("unterminated string at %d", startPos)
}
