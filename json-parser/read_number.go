package jsonparser

import "fmt"

func (l *lexer) readNumber(firstChar byte) (string, error) {
	startingPosition := l.position - 1
	sm := stateMachine{}
	stop, err := sm.next(firstChar)

	if err != nil {
		return "", fmt.Errorf("%w at %d", err, startingPosition)
	}

	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		stop, err = sm.next(char)
		if err != nil {
			return "", fmt.Errorf("%w at %d", err, startingPosition)
		}

		if stop {
			l.position--
			break
		}
	}

	err = sm.end()

	if err != nil {
		return "", fmt.Errorf("%w at %d", err, startingPosition)
	}

	return l.input[startingPosition:l.position], nil
}
