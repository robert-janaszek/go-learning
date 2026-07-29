package course

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("[%s]: %s", v.Field, v.Message)
}

func validate() error {
	return ValidationError{
		Field:   "id",
		Message: "id cannot be 0",
	}
}

func RegisterUser(email, password string) error {
	if !strings.Contains(email, "@") {
		return &ValidationError{Field: "email", Message: "missing @"}
	}

	return nil
}

func badValidate() error {
	var customErr *ValidationError = nil

	return customErr
}
func badValidate1() error {
	// var customErr *ValidationError = nil

	return nil
}

func Day4b() {
	// ex 6
	fmt.Println(validate())

	// ex 7
	err1 := RegisterUser("test", "test")
	if err1 != nil {
		fmt.Println(err1)
	}
	err2 := RegisterUser("test@test.com", "test")
	if err2 != nil {
		fmt.Println(err2)
	}

	// ex 8
	ve, ok := err1.(*ValidationError)
	if ok {
		fmt.Println(ve.Field)
		fmt.Println(ve.Message)
	}

	// ex 9
	err3 := badValidate()
	if err3 != nil {
		fmt.Println(err3) // <nil>
	}

	// ex 10
	err4 := badValidate1()
	if err4 != nil {
		fmt.Println(err4)
	}
}
