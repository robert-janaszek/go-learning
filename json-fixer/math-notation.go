package jsonfixer

import "strings"

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func fixMathNotation(base string, inString bool) string {
	if inString {
		return base
	}

	if strings.HasSuffix(base, ".") && len(base) >= 2 && isDigit(base[len(base)-2]) {
		return base + "0"
	}

	if len(base) >= 3 {
		switch base[len(base)-2:] {
		case "e+", "e-", "E+", "E-":
			if isDigit(base[len(base)-3]) {
				return base + "0"
			}
		}
	}

	if strings.HasSuffix(base, "e") ||
		strings.HasSuffix(base, "E") {
		if len(base) >= 2 && isDigit(base[len(base)-2]) {
			return base + "0"
		}
	}

	if strings.HasSuffix(base, "-") {
		return base + "0"
	}

	return base
}
