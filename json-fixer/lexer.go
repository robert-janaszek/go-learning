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

func (l *Lexer) GetNextSpecialLexem() byte {
	if l.position > len(l.input) {
		return EOF
	}

	for l.position < len(l.input) {
		nextChar := l.input[l.position]

		if nextChar == ParensOpen ||
			nextChar == ParensClose ||
			nextChar == SquareBracketOpen ||
			nextChar == SquareBracketClose ||
			nextChar == StringOpen {
			l.position++
			return nextChar
		}
		l.position++

	}

	return EOF
}
