package jsonparser

func (l *lexer) readKeyword() (token, bool) {
	startingPosition := l.position - 1

	for l.position < len(l.input) {
		char := l.input[l.position]
		l.position++

		if char >= 'a' && char <= 'z' {
			continue
		}

		l.position--
		break
	}

	possibleKeyword := l.input[startingPosition:l.position]

	switch possibleKeyword {
	case "null":
		return token{
			kind: tokenNull,
			lit:  "null",
		}, true
	case "true":
		return token{
			kind: tokenTrue,
			lit:  "true",
		}, true
	case "false":
		return token{
			kind: tokenFalse,
			lit:  "false",
		}, true
	}

	return token{}, false
}
