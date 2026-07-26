package jsonfixer

// Fix appends missing ", }, ] (and a trailing \) so delimiters balance.
// It also trims a trailing comma and completes partial true/false/null outside strings.
// Not always valid JSON.
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

	base := removeTrailingWS(input, lexer.inString)
	base = removeTrailingComma(base, lexer.inString)
	base = fixIncompleteToken(base, lexer.inString)
	return base + suffix, nil
}
