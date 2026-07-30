package jsonparser

import "fmt"

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

func (k tokenKind) String() string {
	switch k {
	case tokenEOF:
		return "EOF"
	case tokenLBrace:
		return "'{'"
	case tokenRBrace:
		return "'}'"
	case tokenLBracket:
		return "'['"
	case tokenRBracket:
		return "']'"
	case tokenColon:
		return "':'"
	case tokenComma:
		return "','"
	case tokenString:
		return "string"
	case tokenNumber:
		return "number"
	case tokenTrue:
		return "true"
	case tokenFalse:
		return "false"
	case tokenNull:
		return "null"
	default:
		return fmt.Sprintf("tokenKind(%d)", k)
	}
}
