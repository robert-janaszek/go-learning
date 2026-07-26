package jsonfixer

import "strings"

func removeTrailingWS(base string, inString bool) string {
	if inString {
		return base
	}

	return strings.TrimRight(base, " \t\n\r")
}
