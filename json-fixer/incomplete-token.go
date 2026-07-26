package jsonfixer

import "strings"

var literals = []string{"true", "false", "null"}

func fixIncompleteToken(base string, inString bool) string {
	if inString {
		return base
	}

	for _, lit := range literals {
		for n := len(lit) - 1; n >= 1; n-- {
			if strings.HasSuffix(base, lit[:n]) {
				return base + lit[n:]
			}
		}
	}
	return base
}
