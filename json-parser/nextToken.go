package jsonparser

func nextToken(l *lexer) (token, error) {
	tok, err := l.next()

	if err != nil {
		return token{}, err
	}

	return tok, nil
}
