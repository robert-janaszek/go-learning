package jsonfixer

import "strings"

func removeTrailingComma(base string, inString bool) string {
	if inString {
		return base
	}

	base = strings.TrimRight(base, " \t\n\r")
	if len(base) > 0 && base[len(base)-1] == ',' {
		return base[:len(base)-1]
	}
	return base
}
