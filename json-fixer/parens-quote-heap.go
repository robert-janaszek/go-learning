package jsonfixer

import "errors"

type ParensQuoteHeap struct {
	heap []byte
}

func (p *ParensQuoteHeap) AddSpecial(char byte) error {
	length := len(p.heap)

	if length == 0 && (char == ParensOpen || char == SquareBracketOpen || char == StringOpen) {
		p.heap = append(p.heap, char)
		return nil
	}

	if length == 0 && (char == ParensClose || char == SquareBracketClose) {
		return errors.New("unexpected char at the beginning, found " + string(char))
	}

	lastChar := p.heap[len(p.heap)-1]
	if lastChar == ParensOpen && char == ParensClose {
		pop(&p.heap)
		return nil
	}

	if lastChar == SquareBracketOpen && char == SquareBracketClose {
		pop(&p.heap)
		return nil
	}

	if char == ParensOpen || char == SquareBracketOpen {
		p.heap = append(p.heap, char)
		return nil
	}

	if char == ParensClose || char == SquareBracketClose {
		return errors.New("unexpected char, found " + string(char))
	}

	return nil
}
