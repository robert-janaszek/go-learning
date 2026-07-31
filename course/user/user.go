package user

import (
	"errors"
	"strings"
)

type User struct {
	email string
}

func (u *User) SetEmail(e string) error {
	if !strings.Contains(e, "@") {
		return errors.New("email has to contain @")
	}

	u.email = e

	return nil
}

func (u User) Email() string {
	return u.email
}
