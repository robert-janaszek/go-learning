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

func (l Lexer) GetNextSpecialLexem() byte {
	return 'a'
}
