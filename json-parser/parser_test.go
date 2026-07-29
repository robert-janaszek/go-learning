package jsonparser

import (
	"reflect"
	"testing"
)

func TestParseScalars(t *testing.T) {
	cases := map[string]any{
		`null`:       nil,
		`true`:       true,
		`false`:      false,
		`"hi"`:       "hi",
		`123`:        float64(123),
		`-1.5`:       float64(-1.5),
		`1e2`:        float64(100),
		`  true  `:   true,
		`"a\"b"`:     `a"b`,
	}

	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			got, err := Parse(in)
			if err != nil {
				t.Fatalf("Parse(%q): unexpected error: %v", in, err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Parse(%q)\n got %#v (%T)\nwant %#v (%T)", in, got, got, want, want)
			}
		})
	}
}

func TestParseArray(t *testing.T) {
	cases := map[string]any{
		`[]`:              []any{},
		`[1]`:             []any{float64(1)},
		`[1, 2]`:          []any{float64(1), float64(2)},
		`[1,2,3]`:         []any{float64(1), float64(2), float64(3)},
		`[true, false]`:   []any{true, false},
		`[null]`:          []any{nil},
		`["a", 1, true]`:  []any{"a", float64(1), true},
		`[[1], [2, 3]]`:   []any{[]any{float64(1)}, []any{float64(2), float64(3)}},
		`[ 1 , 2 ]`:       []any{float64(1), float64(2)},
		`[-1, 1e-2, 3]`:   []any{float64(-1), float64(0.01), float64(3)},
	}

	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			got, err := Parse(in)
			if err != nil {
				t.Fatalf("Parse(%q): unexpected error: %v", in, err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Parse(%q)\n got %#v\nwant %#v", in, got, want)
			}
		})
	}
}

func TestParseArrayErrors(t *testing.T) {
	cases := []string{
		`[`,
		`[1`,
		`[1,`,
		`[1,]`,
		`[1 2]`,
		`[,1]`,
		`[1,,2]`,
		`[1}`,
	}

	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			got, err := Parse(in)
			if err == nil {
				t.Fatalf("Parse(%q): want error, got %#v", in, got)
			}
		})
	}
}
