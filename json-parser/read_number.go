package jsonparser

func (l *lexer) readNumber(firstChar byte) (string, error) {
	startingPosition := l.position - 1
	sm := stateMachine{}
	stop, err := sm.next(firstChar)

	if err != nil {
		return "", err
	}

	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		stop, err = sm.next(char)
		if err != nil {
			return "", err
		}

		if stop {
			l.position--
			break
		}
	}

	err = sm.end()

	if err != nil {
		return "", err
	}

	return l.input[startingPosition:l.position], nil
}
