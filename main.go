package main

import (
	"fmt"

	jsonfixer "github.com/robert-janaszek/go-learning/json-fixer"
)

func main() {
	// day1a()
	// day1b()
	// day1c()
	// day2()
	// day2b()
	// day2c()

	result, _ := jsonfixer.Fix("{{{}[")
	fmt.Println(result)
}
