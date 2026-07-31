package jsonparser

import "fmt"

func expectToken(tok token, l *lexer, exp tokenKind) error {
	if tok.kind != exp {
		return fmt.Errorf("expected token %s, found: %q at %d", exp, tok.lit, tok.pos)
	}

	return nil
}
