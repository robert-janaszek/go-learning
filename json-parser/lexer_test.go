package jsonparser

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	cases := map[string][]token{
		``:    {},
		`{}`:  {{tokenLBrace, "{"}, {tokenRBrace, "}"}},
		`[]`:  {{tokenLBracket, "["}, {tokenRBracket, "]"}},
		`:`:   {{tokenColon, ":"}},
		`,`:   {{tokenComma, ","}},
		`true`:  {{tokenTrue, "true"}},
		`false`: {{tokenFalse, "false"}},
		`null`:  {{tokenNull, "null"}},
		`123`:    {{tokenNumber, "123"}},
		`-1`:     {{tokenNumber, "-1"}},
		`-1.5`:   {{tokenNumber, "-1.5"}},
		`0.5`:    {{tokenNumber, "0.5"}},
		`1e2`:    {{tokenNumber, "1e2"}},
		`1E2`:    {{tokenNumber, "1E2"}},
		`1e+2`:   {{tokenNumber, "1e+2"}},
		`1e-2`:   {{tokenNumber, "1e-2"}},
		`1.5e-3`: {{tokenNumber, "1.5e-3"}},
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
		`"hi"`: {{tokenString, "hi"}},
		`"a\"b"`: {{tokenString, `a"b`}},
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
		got := collectTokens(in)
		if len(got) == 0 && len(want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("lexer(%q)\n got %#v\nwant %#v", in, got, want)
		}
	}
}

func collectTokens(input string) []token {
	l := lexer{}
	l.start(input)
	var out []token
	for {
		tok, ok := l.next()
		if !ok {
			return out
		}
		out = append(out, tok)
	}
}
