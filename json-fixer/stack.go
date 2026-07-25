package jsonfixer

import "errors"

type Stack struct {
	stack []byte
}

func (p *Stack) AddSpecial(char byte) error {
	length := len(p.stack)

	if length == 0 && (char == braceOpen || char == squareBracketOpen || char == stringOpen) {
		p.stack = append(p.stack, char)
		return nil
	}

	if length == 0 && (char == braceClose || char == squareBracketClose) {
		return errors.New("unexpected char at the beginning, found " + string(char))
	}

	if length == 0 {
		return nil
	}

	lastChar := p.stack[len(p.stack)-1]
	if lastChar == braceOpen && char == braceClose {
		pop(&p.stack)
		return nil
	}

	if lastChar == squareBracketOpen && char == squareBracketClose {
		pop(&p.stack)
		return nil
	}

	if char == braceOpen || char == squareBracketOpen {
		p.stack = append(p.stack, char)
		return nil
	}

	if char == braceClose || char == squareBracketClose {
		return errors.New("unexpected char, found " + string(char))
	}

	return nil
}
