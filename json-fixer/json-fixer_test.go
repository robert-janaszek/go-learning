package jsonfixer

import "testing"

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
