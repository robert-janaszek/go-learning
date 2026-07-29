package course

import (
	"errors"
	"fmt"
	"os"
)

func ReadConfig(path string) error {
	_, err := os.ReadFile(path)

	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	return nil
}

func validateForm() error {
	err1 := fmt.Errorf("field ID does not exist")
	err2 := fmt.Errorf("field Name does not exist")
	err3 := fmt.Errorf("field Path does not exist")

	errs := []error{err1, err2, err3}

	return errors.Join(errs...)
}

func Repository() error {
	return ErrNotFound
}

func Service() error {
	err := Repository()
	return fmt.Errorf("user service: %w", err)
}

func Handler() int {
	err := Service()

	if errors.Is(err, ErrNotFound) {
		return 404
	}

	return 200
}

func Day4c() {
	// ex 11
	err := ReadConfig("non-existing")
	if err != nil {
		fmt.Println(err)
	}

	// ex 12
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println(err)
	}

	// ex 13
	var pathErr *os.PathError
	ok := errors.As(err, &pathErr)
	if ok {
		fmt.Println(pathErr.Path)
	}

	// ex 14
	fmt.Println(validateForm())

	// ex 15
	fmt.Println(Handler())
}
