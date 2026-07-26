package jsonfixer

func removeTrailingComma(base string, inString bool) string {
	if inString {
		return base
	}

	if len(base) > 0 && base[len(base)-1] == ',' {
		return base[:len(base)-1]
	}
	return base
}
