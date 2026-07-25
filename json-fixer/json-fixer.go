package jsonfixer

func Fix(input string) (string, error) {
	stk := stack{}
	var err error

	lexer := lexer{}
	lexer.start(input)

	for {
		lex, ok := lexer.next()

		if !ok {
			break
		}

		err = stk.addSpecial(lex)

		if err != nil {
			return "", err
		}
	}
	// fmt.Printf("%s\n", stk.stack)

	suffix := stk.close()

	return input + suffix, nil
}
