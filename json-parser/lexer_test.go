package jsonparser

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	cases := map[string][]token{
		``:        {},
		`{}`:      {{tokenLBrace, "{"}, {tokenRBrace, "}"}},
		`[]`:      {{tokenLBracket, "["}, {tokenRBracket, "]"}},
		`:`:       {{tokenColon, ":"}},
		`,`:       {{tokenComma, ","}},
		`true`:    {{tokenTrue, "true"}},
		`false`:   {{tokenFalse, "false"}},
		`null`:    {{tokenNull, "null"}},
		`123`:     {{tokenNumber, "123"}},
		`-1`:      {{tokenNumber, "-1"}},
		`-1.5`:    {{tokenNumber, "-1.5"}},
		`0.5`:     {{tokenNumber, "0.5"}},
		`1e2`:     {{tokenNumber, "1e2"}},
		`1E2`:     {{tokenNumber, "1E2"}},
		`1e+2`:    {{tokenNumber, "1e+2"}},
		`1e-2`:    {{tokenNumber, "1e-2"}},
		`1.5e-3`:  {{tokenNumber, "1.5e-3"}},
		`-2.5E+1`: {{tokenNumber, "-2.5E+1"}},
		`[1, 2]`: {
			{tokenLBracket, "["},
			{tokenNumber, "1"},
			{tokenComma, ","},
			{tokenNumber, "2"},
			{tokenRBracket, "]"},
		},
		`[-1, 1e-2, 3]`: {
			{tokenLBracket, "["},
			{tokenNumber, "-1"},
			{tokenComma, ","},
			{tokenNumber, "1e-2"},
			{tokenComma, ","},
			{tokenNumber, "3"},
			{tokenRBracket, "]"},
		},
		`{"n":-4.2e+1}`: {
			{tokenLBrace, "{"},
			{tokenString, "n"},
			{tokenColon, ":"},
			{tokenNumber, "-4.2e+1"},
			{tokenRBrace, "}"},
		},
		`"hi"`:     {{tokenString, "hi"}},
		`"a\"b"`:   {{tokenString, `a"b`}},
		`"a\\b"`:   {{tokenString, `a\b`}},
		`"\\"`:     {{tokenString, `\`}},
		`"\\\""`:   {{tokenString, `\"`}},
		`"a\\\"b"`: {{tokenString, `a\"b`}},
		`[1, true, "a"]`: {
			{tokenLBracket, "["},
			{tokenNumber, "1"},
			{tokenComma, ","},
			{tokenTrue, "true"},
			{tokenComma, ","},
			{tokenString, "a"},
			{tokenRBracket, "]"},
		},
		`{"a":1}`: {
			{tokenLBrace, "{"},
			{tokenString, "a"},
			{tokenColon, ":"},
			{tokenNumber, "1"},
			{tokenRBrace, "}"},
		},
		`  { "x" : null }  `: {
			{tokenLBrace, "{"},
			{tokenString, "x"},
			{tokenColon, ":"},
			{tokenNull, "null"},
			{tokenRBrace, "}"},
		},
	}

	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			got, err := collectTokens(in)
			if err != nil {
				t.Fatalf("lexer(%q): unexpected error: %v", in, err)
			}
			if len(got) == 0 && len(want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("lexer(%q)\n got %#v\nwant %#v", in, got, want)
			}
		})
	}
}

func TestLexerEOF(t *testing.T) {
	l := lexer{}
	l.start("")

	tok, err := l.next()
	if err != nil {
		t.Fatalf("empty input: unexpected error: %v", err)
	}
	if tok.kind != tokenEOF {
		t.Fatalf("empty input: got kind %v, want tokenEOF", tok.kind)
	}

	l.start("   \t\n")
	tok, err = l.next()
	if err != nil {
		t.Fatalf("whitespace-only: unexpected error: %v", err)
	}
	if tok.kind != tokenEOF {
		t.Fatalf("whitespace-only: got kind %v, want tokenEOF", tok.kind)
	}
}

func TestLexerErrors(t *testing.T) {
	cases := []string{
		`@`,
		`#`,
		`"`,
		`"abc`,
		`"\`,
		`tru`,
		`tree`,
		`True`,
		`1e2e3`,
	}

	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			_, err := collectTokens(in)
			if err == nil {
				t.Fatalf("lexer(%q): want error", in)
			}
		})
	}
}

func TestLexerInvalidNumbers(t *testing.T) {
	cases := []string{
		`-`,
		`01`,
		`-01`,
		`00`,
		`1.`,
		`0.`,
		`-1.`,
		`.5`,
		`1e`,
		`1E`,
		`1e+`,
		`1e-`,
		`-1e`,
		`1.e2`,
		`1e.2`,
		`0e`,
		`0e+`,
	}

	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			_, err := collectTokens(in)
			if err == nil {
				t.Fatalf("lexer(%q): want error for invalid number", in)
			}
		})
	}
}

// collectTokens returns all tokens until EOF.
// Lexical errors are returned as err (not silently treated as EOF).
func collectTokens(input string) ([]token, error) {
	l := lexer{}
	l.start(input)
	var out []token
	for {
		tok, err := l.next()
		if err != nil {
			return out, err
		}
		if tok.kind == tokenEOF {
			return out, nil
		}
		out = append(out, tok)
	}
}
