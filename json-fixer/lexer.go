package jsonfixer

type lexer struct {
	position int
	input    string
}

func (l *lexer) start(input string) {
	l.input = input
	l.position = 0
}

func (l *lexer) next() (byte, bool) {
	for l.position < len(l.input) {
		nextChar := l.input[l.position]
		l.position++

		switch nextChar {
		case braceOpen, braceClose, squareBracketOpen, squareBracketClose, stringOpen:
			return nextChar, true
		}
	}

	return 0, false
}
