package jsonfixer

// Fix appends missing ", }, ] (and a trailing \) so delimiters balance.
// It does not produce fully valid JSON in all cases.
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

	suffix := stk.close()

	return input + suffix, nil
}
