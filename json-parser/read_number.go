package jsonparser

import "errors"

// TODO: add more defences for numbers, RFC
func (l *lexer) readNumber() (string, error) {
	eConsumed := false
	dotConsumed := false
	startingPosition := l.position - 1

	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		if char >= '0' && char <= '9' {
			continue
		}

		switch char {
		case 'e', 'E':
			if eConsumed {
				return "", errors.New("incorrect number, found 2 'e'")
			}

			eConsumed = true

			if l.position < len(l.input) {
				nextChar := l.input[l.position]

				switch nextChar {
				case '-', '+':
					l.position++
				}
			}

			continue
		}

		if char == '.' && !dotConsumed {
			dotConsumed = true
			continue
		}

		l.position--
		break
	}

	return l.input[startingPosition:l.position], nil
}
