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

	if char == backslash {
		s.push(char)
		return nil
	}

	if length == 0 {
		switch char {
		case braceOpen, squareBracketOpen:
			s.push(char)
			return nil

		case braceClose, squareBracketClose:
			return fmt.Errorf("unexpected %q", string(char))

		case stringOpen:
			s.push(char)
			return nil
		}

		return nil
	}

	lastChar := s.stack[len(s.stack)-1]

	if lastChar == braceOpen && char == braceClose {
		s.pop()
		return nil
	}

	if lastChar == squareBracketOpen && char == squareBracketClose {
		s.pop()
		return nil
	}

	switch char {
	case stringOpen:
		if s.stack[len(s.stack)-1] == stringOpen {
			s.pop()
			return nil
		}
		s.push(char)
		return nil
	case braceOpen, squareBracketOpen:
		s.push(char)
		return nil

	case braceClose, squareBracketClose:
		return fmt.Errorf("unexpected %q", string(char))

	}

	return nil
}

func (s *stack) close() string {
	var enclosure strings.Builder

	for len(s.stack) > 0 {
		switch s.pop() {
		case braceOpen:
			enclosure.WriteString(string(braceClose))
		case squareBracketOpen:
			enclosure.WriteString(string(squareBracketClose))
		case stringOpen:
			enclosure.WriteString("\"")
		case backslash:
			enclosure.WriteString("\\")
		}
	}
	return enclosure.String()
}
