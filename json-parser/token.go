package jsonparser

type tokenKind int

const (
	tokenEOF tokenKind = iota
	tokenLBrace
	tokenRBrace
	tokenLBracket
	tokenRBracket
	tokenColon
	tokenComma
	tokenString
	tokenNumber
	tokenTrue
	tokenFalse
	tokenNull
)

type token struct {
	kind tokenKind
	lit  string
}

func tok(kind tokenKind, lit string) token {
	return token{
		kind: kind,
		lit:  lit,
	}
}
