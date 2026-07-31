package jsonparser

import "fmt"

func expectToken(tok token, exp tokenKind) error {
	if tok.kind != exp {
		return fmt.Errorf("expected %s, found %q at %d", exp, tok.lit, tok.pos)
	}

	return nil
}
