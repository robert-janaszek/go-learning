package days

import (
	"errors"
	"fmt"
	"strconv"
)

func ValidateAge(age int) error {
	if age < 0 || age > 120 {
		return fmt.Errorf("age %d is out of range [0-120]", age)
	}

	return nil
}

func ProcessOrder(id int, amount float64) error {
	if id <= 0 {
		return fmt.Errorf("id %d cannot be 0 or lower", id)
	}

	if amount <= 0 {
		return fmt.Errorf("amount %v cannot be 0 or lower", amount)
	}

	if amount >= 10000 {
		return fmt.Errorf("amount %v cannot be 10000 or higher", amount)
	}

	return nil
}

var ErrNotFound = errors.New("item not found")
var ErrPermissionDenied = errors.New("permission denied")

func FindUser(id int) (*user, error) {
	if id == 0 {
		return nil, ErrNotFound
	}

	return nil, nil
}

func Day5a() {
	// ex 1, 2
	ages := []int{10, -1, 121}

	for _, age := range ages {
		err := ValidateAge(age)
		if err != nil {
			fmt.Println(err)
		}
	}

	// ex 3

	err1 := ProcessOrder(0, 1000)
	if err1 != nil {
		fmt.Println(err1)
	}
	err2 := ProcessOrder(1, 0)
	if err2 != nil {
		fmt.Println(err2)
	}
	err3 := ProcessOrder(2, 1000)
	if err3 != nil {
		fmt.Println(err3)
	}

	// ex4

	_, err4 := FindUser(0)

	if err4 == ErrNotFound {
		fmt.Println(err4)
	}

	// ex5

	val, _ := strconv.Atoi("123")
	fmt.Println(val)
}
