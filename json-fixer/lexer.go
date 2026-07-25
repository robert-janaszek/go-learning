package jsonfixer

type Lexer struct {
	position int
	input    string
}

const EOF byte = 'E'

func (l *Lexer) Start(input string) {
	l.input = input
	l.position = 0
}

func (l *Lexer) Next() byte {
	if l.position >= len(l.input) {
		return EOF
	}

	for l.position < len(l.input) {
		nextChar := l.input[l.position]

		if nextChar == braceOpen ||
			nextChar == braceClose ||
			nextChar == squareBracketOpen ||
			nextChar == squareBracketClose ||
			nextChar == stringOpen {
			l.position++
			return nextChar
		}
		l.position++

	}

	return EOF
}
