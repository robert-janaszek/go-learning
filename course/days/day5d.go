package days

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func MustParseURL(rawURL string) {
	if rawURL == "" {
		panic("invalid URL")
	}
}

func SafeExecute(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered: %s\n", r)
		}
	}()
	fn()
}

func SafeExecuteIntoErr(fn func()) (err error) {
	// var err error = nil
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	fn()

	return err
}

func openFile() error {
	file, err := os.Open("./course/days/day5d.go")
	if err != nil {
		fmt.Println("not opened")
		return err
	}
	fmt.Println("opened")
	defer file.Close()

	if true {
		return errors.New("unexpected error")
	}

	return nil
}

func ExecuteWithLogging(fn func() error) error {
	err := fn()
	if err != nil {
		log.Println(err)
	}

	return err
}

func Day5d() {
	// ex 16
	MustParseURL("test.com")
	// MustParseURL("")

	// ex 17
	SafeExecute(func() { MustParseURL("") })

	// ex 18
	err := SafeExecuteIntoErr(func() { MustParseURL("") })
	if err != nil {
		fmt.Println(err)
	}

	// ex 19
	openFile()

	// ex 20
	err = ExecuteWithLogging(func() error {
		return SafeExecuteIntoErr(func() { MustParseURL("") })
	})

	if err != nil {
		// log
	}
}
