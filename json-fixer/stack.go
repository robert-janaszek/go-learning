package jsonfixer

import (
	"errors"
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

func (s *stack) append(item byte) {
	s.stack = append(s.stack, item)
}

func (s *stack) addSpecial(char byte) error {
	length := len(s.stack)

	if char == backslash {
		s.append(char)
		return nil
	}

	if length == 0 {
		switch char {
		case braceOpen, squareBracketOpen:
			s.append(char)
			return nil

		case braceClose, squareBracketClose:
			return errors.New("unexpected char at the beginning, found " + string(char))

		case stringOpen:
			s.append(char)
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
		s.append(char)
		return nil
	case braceOpen, squareBracketOpen:
		s.append(char)
		return nil

	case braceClose, squareBracketClose:
		return errors.New("unexpected char, found " + string(char))

	}

	return nil
}

func (s *stack) close() string {
	var enclosure string
	for len(s.stack) > 0 {
		switch s.pop() {
		case braceOpen:
			enclosure += string(braceClose)
		case squareBracketOpen:
			enclosure += string(squareBracketClose)
		case stringOpen:
			enclosure += "\""
		case backslash:
			enclosure += "\\"
		}
	}
	return enclosure
}
