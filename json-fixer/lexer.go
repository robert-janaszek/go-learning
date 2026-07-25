package jsonfixer

type lexer struct {
	position int
	inString bool
	input    string
}

func (l *lexer) start(input string) {
	l.input = input
	l.inString = false
	l.position = 0
}

func (l *lexer) next() (byte, bool) {
	for l.position < len(l.input) {
		nextChar := l.input[l.position]
		l.position++

		if l.inString && nextChar == backslash {
			if l.position >= len(l.input) {
				return backslash, true
			}
			l.position++
			continue
		}

		switch nextChar {
		case stringOpen:
			l.inString = !l.inString
			return stringOpen, true
		case braceOpen, braceClose, squareBracketOpen, squareBracketClose:
			return nextChar, true
		}
	}

	return 0, false
}
