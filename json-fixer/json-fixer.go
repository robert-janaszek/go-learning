package jsonfixer

// Fix appends missing ", }, ] (and a trailing \) so delimiters balance.
// Outside strings it also trims a trailing comma, completes partial
// true/false/null, and repairs truncated numbers (1., 1e, 1e+, -).
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
	base = fixMathNotation(base, lexer.inString)
	return base + suffix, nil
}
