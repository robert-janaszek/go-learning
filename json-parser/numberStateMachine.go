package jsonparser

import (
	"errors"
	"fmt"
)

type stateMachine struct {
	currentState stateMachineState
}

type stateMachineState int

const (
	start stateMachineState = iota
	minus
	int_zero
	integer
	dot
	frac
	exp
	exp_sign
	exp_digits
)

func (s *stateMachine) startOnChar(char byte) error {
	switch char {
	case '-':
		s.currentState = minus
		return nil
	case '0':
		s.currentState = int_zero
		return nil
	}

	if char >= '1' && char <= '9' {
		s.currentState = integer
		return nil
	}

	return fmt.Errorf("incorrect token at start, found %q", char)
}

func (s *stateMachine) minusOnChar(char byte) error {
	if char == '0' {
		s.currentState = int_zero
		return nil
	}
	if char >= '1' && char <= '9' {
		s.currentState = integer
		return nil
	}

	return fmt.Errorf("incorrect token after '-', found %q", char)
}

func (s *stateMachine) intZeroOnChar(char byte) (bool, error) {
	switch char {
	case '.':
		s.currentState = dot
		return false, nil
	case 'e', 'E':
		s.currentState = exp
		return false, nil
	}

	if char >= '0' && char <= '9' {
		return false, fmt.Errorf("unexpected digit after zero, found: %q", char)
	}

	return true, nil
}

func (s *stateMachine) intOnChar(char byte) (bool, error) {
	switch char {
	case '.':
		s.currentState = dot
		return false, nil
	case 'e', 'E':
		s.currentState = exp
		return false, nil
	}

	if char >= '0' && char <= '9' {
		// no change
		return false, nil
	}

	return true, nil
}

func (s *stateMachine) dotOnChar(char byte) (bool, error) {
	if char >= '0' && char <= '9' {
		s.currentState = frac
		return false, nil
	}

	return false, fmt.Errorf("unexpected char after dot, found %q", char)
}

func (s *stateMachine) fracOnChar(char byte) (bool, error) {
	if char >= '0' && char <= '9' {
		return false, nil
	}

	if char == 'e' || char == 'E' {
		s.currentState = exp
		return false, nil
	}

	return true, nil
}

func (s *stateMachine) expOnChar(char byte) (bool, error) {
	if char == '-' || char == '+' {
		s.currentState = exp_sign
		return false, nil
	}

	if char >= '0' && char <= '9' {
		s.currentState = exp_digits
		return false, nil
	}

	return false, fmt.Errorf("incorrect number, found: %q", char)
}

func (s *stateMachine) expSignOnChar(char byte) (bool, error) {
	if char >= '0' && char <= '9' {
		s.currentState = exp_digits
		return false, nil
	}

	return false, fmt.Errorf("incorrect number, found %q", char)
}
func (s *stateMachine) expDigitOnChar(char byte) (bool, error) {
	if char >= '0' && char <= '9' {
		return false, nil
	}

	return true, nil
}

func (s *stateMachine) next(char byte) (bool, error) {
	var err error
	var stop bool
	switch s.currentState {
	case start:
		err = s.startOnChar(char)
		return false, err
	case minus:
		err = s.minusOnChar(char)
		return false, err
	case int_zero:
		stop, err = s.intZeroOnChar(char)
		return stop, err
	case integer:
		stop, err = s.intOnChar(char)
		return stop, err
	case dot:
		stop, err = s.dotOnChar(char)
		return stop, err
	case frac:
		stop, err = s.fracOnChar(char)
		return stop, err
	case exp:
		stop, err = s.expOnChar(char)
		return stop, err
	case exp_sign:
		stop, err = s.expSignOnChar(char)
		return stop, err
	case exp_digits:
		stop, err = s.expDigitOnChar(char)
		return stop, err
	}

	return false, errors.New("invalid state")
}

func (s *stateMachine) end() error {
	switch s.currentState {
	case start:
		return errors.New("incorrect number")
	case minus:
		return errors.New("incorrect number, no digit after '-'")
	case dot:
		return errors.New("incorrect number, no digit after '.'")
	case exp:
		return errors.New("incorrect number, incorrect mathematical expression")
	case exp_sign:
		return errors.New("incorrect number, incorrect mathematical expression, no digit after sign")
	}

	return nil
}
