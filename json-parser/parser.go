package jsonparser

func Parse(input string) (any, error) {
	l := lexer{}
	l.start(input)
	for {
		_, ok := l.next()

		if !ok {
			break
		}
	}
	return nil, nil
}
