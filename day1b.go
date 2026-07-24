package main

import (
	"fmt"
	"strings"
)

func swap(a, b *int) {
	temp := *a
	*a = *b
	*b = temp
}

func betterSwap(a, b int) (int, int) {
	return b, a
}

func uppercase(str string) string {
	return strings.ToUpper(str)
}

func uppercaseMutate(str *string) {
	*str = strings.ToUpper(*str)
}

func day1b() {
	// ex 11
	a := 1
	b := 2

	swap(&a, &b)
	fmt.Println(a, b)

	a = 1
	b = 2
	a, b = betterSwap(a, b)
	fmt.Println(a, b)

	// ex 12

	str := "aBcDe"
	newStr := uppercase(str)
	uppercaseMutate(&str)

	fmt.Println(str, newStr)
}
