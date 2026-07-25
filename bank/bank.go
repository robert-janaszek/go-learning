package bank

import "errors"

type Account struct {
	balance float64
}

func (a *Account) Deposit(amount float64) error {
	if amount < 0 {
		return errors.New("cannot deposit negative amount")
	}
	a.balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount < 0 {
		return errors.New("cannot withdraw negative amount")
	}

	if a.balance-amount < 0 {
		return errors.New("insufficient funds")
	}

	a.balance -= amount

	return nil
}

func (a Account) Balance() float64 {
	return a.balance
}
