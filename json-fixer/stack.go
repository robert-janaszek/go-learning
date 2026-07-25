package jsonfixer

import (
	"fmt"
	"strings"
)

type stack struct {
	stack []byte
}

func (s *stack) pop() byte {
	length := len(s.stack)
	lastChar := s.stack[length-1]
	s.stack = s.stack[:length-1]

	return lastChar
}

func (s *stack) push(item byte) {
	s.stack = append(s.stack, item)
}

func (s *stack) addSpecial(char byte) error {
	length := len(s.stack)
	var lastChar byte
	if length > 0 {
		lastChar = s.stack[len(s.stack)-1]
	}

	switch char {
	case backslash:
		s.push(char)
		return nil
	case braceOpen, squareBracketOpen:
		s.push(char)
		return nil
	case stringOpen:
		if lastChar == stringOpen {
			s.pop()
			return nil
		}
		s.push(char)
		return nil

	case braceClose:
		if lastChar == braceOpen {
			s.pop()
			return nil
		}

		return fmt.Errorf("unexpected %q", char)
	case squareBracketClose:

		if lastChar == squareBracketOpen {
			s.pop()
			return nil
		}
		return fmt.Errorf("unexpected %q", char)
	}

	return nil
}

func (s *stack) close() string {
	var enclosure strings.Builder

	for len(s.stack) > 0 {
		switch s.pop() {
		case braceOpen:
			enclosure.WriteByte(braceClose)
		case squareBracketOpen:
			enclosure.WriteByte(squareBracketClose)
		case stringOpen:
			enclosure.WriteByte('"')
		case backslash:
			enclosure.WriteByte('\\')
		}
	}
	return enclosure.String()
}
