package jsonfixer

import "testing"

func TestFixTrailingComma(t *testing.T) {
	cases := map[string]string{
		`[1,`:     `[1]`,
		`{"a":1,`: `{"a":1}`,
		`[1, `:    `[1]`,
		`"hello,`: `"hello,"`,
	}
	for in, want := range cases {
		got, err := Fix(in)
		if err != nil || got != want {
			t.Errorf("Fix(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
}

func TestFixIncompleteToken(t *testing.T) {
	cases := map[string]string{
		`n`:         `null`,
		`nu`:        `null`,
		`nul`:       `null`,
		`null`:      `null`,
		`t`:         `true`,
		`tr`:        `true`,
		`tru`:       `true`,
		`true`:      `true`,
		`f`:         `false`,
		`fa`:        `false`,
		`fal`:       `false`,
		`fals`:      `false`,
		`false`:     `false`,
		`{"a":tru`:  `{"a":true}`,
		`{"a":nul`:  `{"a":null}`,
		`[fals`:     `[false]`,
		`{"a":tru `: `{"a":true}`,
		`"tru`:      `"tru"`,
		`{"a":"nul`: `{"a":"nul"}`,
	}
	for in, want := range cases {
		got, err := Fix(in)
		if err != nil || got != want {
			t.Errorf("Fix(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
}

func TestFix(t *testing.T) {
	cases := map[string]string{
		"{}":     "{}",
		"[{}]":   "[{}]",
		"{":      "{}",
		"{{":     "{{}}",
		"[":      "[]",
		"{{{}[":  "{{{}[]}}",
		`{1,2`:   `{1,2}`,
		`{"a`:    `{"a"}`,
		`{"a"`:   `{"a"}`,
		`{"a"}`:  `{"a"}`,
		`"hi`:    `"hi"`,
		`{"x":{`: `{"x":{}}`,
		`{"a\":`: `{"a\":"}`,
		`"abc\`:  `"abc\\"`,
		`"a\"b`:  `"a\"b"`,
		`"\`:     `"\\"`,
		`{"a\\`:  `{"a\\"}`,
	}

	for in, want := range cases {
		got, err := Fix(in)
		if err != nil || got != want {
			t.Errorf("Fix(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
}

func TestFixIncompleteNumber(t *testing.T) {
	cases := map[string]string{
		`1.`:       `1.0`,
		`0.`:       `0.0`,
		`-1.`:      `-1.0`,
		`{"a":1.`:  `{"a":1.0}`,
		`[1.`:      `[1.0]`,
		`1.0`:      `1.0`,
		`1.5`:      `1.5`,
		`1e2`:      `1e2`,
		`1e+2`:     `1e+2`,
		`1e`:       `1e0`,
		`1E`:       `1E0`,
		`1e+`:      `1e+0`,
		`1e-`:      `1e-0`,
		`1.2e`:     `1.2e0`,
		`{"a":1e`:  `{"a":1e0}`,
		`[1e+`:     `[1e+0]`,
		`-`:        `-0`,
		`{"a":-`:   `{"a":-0}`,
		`[-`:       `[-0]`,
		`"1.`:      `"1."`,
		`{"a":"1e`: `{"a":"1e"}`,
		`1. `:      `1.0`,
		`1e+ `:     `1e+0`,
	}
	for in, want := range cases {
		got, err := Fix(in)
		if err != nil || got != want {
			t.Errorf("Fix(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
}

func TestFixError(t *testing.T) {
	for _, in := range []string{"}", "{]", "{[}"} {
		if _, err := Fix(in); err == nil {
			t.Errorf("Fix(%q): want error", in)
		}
	}
}
