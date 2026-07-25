package jsonfixer

import (
	"errors"
)

type stack struct {
	stack []byte
}

func (p *stack) pop() byte {
	length := len(p.stack)
	lastChar := p.stack[length-1]
	p.stack = p.stack[:length-1]

	return lastChar
}

func (p *stack) append(item byte) {
	p.stack = append(p.stack, item)
}

func (p *stack) addSpecial(char byte) error {
	length := len(p.stack)

	if length == 0 {
		switch char {
		case braceOpen, squareBracketOpen, stringOpen:
			p.append(char)
			return nil
		case braceClose, squareBracketClose:
			return errors.New("unexpected char at the beginning, found " + string(char))
		}

		return nil
	}

	lastChar := p.stack[len(p.stack)-1]
	if lastChar == braceOpen && char == braceClose {
		p.pop()
		return nil
	}

	if lastChar == squareBracketOpen && char == squareBracketClose {
		p.pop()
		return nil
	}

	switch char {
	case braceOpen, squareBracketOpen:
		p.append(char)
		return nil

	case braceClose, squareBracketClose:
		return errors.New("unexpected char, found " + string(char))
	}

	return nil
}

func (p *stack) close() string {
	var enclosure string
	var popChar byte
	for len(p.stack) > 0 {
		popChar = p.pop()

		switch popChar {
		case braceOpen:
			enclosure += string(braceClose)
		case squareBracketOpen:
			enclosure += string(squareBracketClose)
		}
	}
	return enclosure
}
