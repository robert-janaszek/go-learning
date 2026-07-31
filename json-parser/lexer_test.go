package jsonparser

import (
	"testing"
)

func TestLexer(t *testing.T) {
	cases := map[string][]token{
		``:        {},
		`{}`:      {{kind: tokenLBrace, lit: "{"}, {kind: tokenRBrace, lit: "}"}},
		`[]`:      {{kind: tokenLBracket, lit: "["}, {kind: tokenRBracket, lit: "]"}},
		`:`:       {{kind: tokenColon, lit: ":"}},
		`,`:       {{kind: tokenComma, lit: ","}},
		`true`:    {{kind: tokenTrue, lit: "true"}},
		`false`:   {{kind: tokenFalse, lit: "false"}},
		`null`:    {{kind: tokenNull, lit: "null"}},
		`123`:     {{kind: tokenNumber, lit: "123"}},
		`-1`:      {{kind: tokenNumber, lit: "-1"}},
		`-1.5`:    {{kind: tokenNumber, lit: "-1.5"}},
		`0.5`:     {{kind: tokenNumber, lit: "0.5"}},
		`1e2`:     {{kind: tokenNumber, lit: "1e2"}},
		`1E2`:     {{kind: tokenNumber, lit: "1E2"}},
		`1e+2`:    {{kind: tokenNumber, lit: "1e+2"}},
		`1e-2`:    {{kind: tokenNumber, lit: "1e-2"}},
		`1.5e-3`:  {{kind: tokenNumber, lit: "1.5e-3"}},
		`-2.5E+1`: {{kind: tokenNumber, lit: "-2.5E+1"}},
		`[1, 2]`: {
			{kind: tokenLBracket, lit: "["},
			{kind: tokenNumber, lit: "1"},
			{kind: tokenComma, lit: ","},
			{kind: tokenNumber, lit: "2"},
			{kind: tokenRBracket, lit: "]"},
		},
		`[-1, 1e-2, 3]`: {
			{kind: tokenLBracket, lit: "["},
			{kind: tokenNumber, lit: "-1"},
			{kind: tokenComma, lit: ","},
			{kind: tokenNumber, lit: "1e-2"},
			{kind: tokenComma, lit: ","},
			{kind: tokenNumber, lit: "3"},
			{kind: tokenRBracket, lit: "]"},
		},
		`{"n":-4.2e+1}`: {
			{kind: tokenLBrace, lit: "{"},
			{kind: tokenString, lit: "n"},
			{kind: tokenColon, lit: ":"},
			{kind: tokenNumber, lit: "-4.2e+1"},
			{kind: tokenRBrace, lit: "}"},
		},
		`"hi"`:     {{kind: tokenString, lit: "hi"}},
		`"a\"b"`:   {{kind: tokenString, lit: `a"b`}},
		`"a\\b"`:   {{kind: tokenString, lit: `a\b`}},
		`"\\"`:     {{kind: tokenString, lit: `\`}},
		`"\\\""`:   {{kind: tokenString, lit: `\"`}},
		`"a\\\"b"`: {{kind: tokenString, lit: `a\"b`}},
		`[1, true, "a"]`: {
			{kind: tokenLBracket, lit: "["},
			{kind: tokenNumber, lit: "1"},
			{kind: tokenComma, lit: ","},
			{kind: tokenTrue, lit: "true"},
			{kind: tokenComma, lit: ","},
			{kind: tokenString, lit: "a"},
			{kind: tokenRBracket, lit: "]"},
		},
		`{"a":1}`: {
			{kind: tokenLBrace, lit: "{"},
			{kind: tokenString, lit: "a"},
			{kind: tokenColon, lit: ":"},
			{kind: tokenNumber, lit: "1"},
			{kind: tokenRBrace, lit: "}"},
		},
		`  { "x" : null }  `: {
			{kind: tokenLBrace, lit: "{"},
			{kind: tokenString, lit: "x"},
			{kind: tokenColon, lit: ":"},
			{kind: tokenNull, lit: "null"},
			{kind: tokenRBrace, lit: "}"},
		},
	}

	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			got, err := collectTokens(in)
			if err != nil {
				t.Fatalf("lexer(%q): unexpected error: %v", in, err)
			}
			if !tokensKindLitEqual(got, want) {
				t.Errorf("lexer(%q)\n got %#v\nwant %#v", in, got, want)
			}
		})
	}
}

func TestLexerPos(t *testing.T) {
	got, err := collectTokens(`  {"a":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantPos := []int{2, 3, 6, 7, 8} // { "a" : 1 }
	if len(got) != len(wantPos) {
		t.Fatalf("got %d tokens, want %d: %#v", len(got), len(wantPos), got)
	}
	for i, pos := range wantPos {
		if got[i].pos != pos {
			t.Errorf("token[%d] (%s %q): pos=%d, want %d", i, got[i].kind, got[i].lit, got[i].pos, pos)
		}
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
	if tok.pos != 0 {
		t.Errorf("empty input EOF pos=%d, want 0", tok.pos)
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

func tokensKindLitEqual(got, want []token) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i].kind != want[i].kind || got[i].lit != want[i].lit {
			return false
		}
	}
	return true
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
