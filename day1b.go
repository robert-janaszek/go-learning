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

func safeDivide(a, b int, result *float64) bool {
	if b == 0 {
		return false
	}

	*result = float64(a) / float64(b)

	return true
}

func safeDivide2(a, b int) (float64, bool) {
	if b == 0 {
		return 0, false
	}

	return float64(a) / float64(b), true
}

func increment(a *int) {
	*a++
}

func allocateInt(a int) *int {
	localVal := a

	return &localVal
}

func localInt(a int) int {
	localVal := a

	return localVal
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

	// ex 13

	a1 := 12
	b1 := 6
	var result float64
	c1 := 0

	done := safeDivide(a1, b1, &result)

	fmt.Println(done, result) // true 2

	done = safeDivide(a1, c1, &result)

	fmt.Println(done, result) // false 2

	result, done = safeDivide2(a1, 3)

	fmt.Println(done, result)

	// ex 14
	a2 := 0

	for i := 0; i < 3; i++ {
		increment(&a2)
	}

	fmt.Println(a2)

	a2 = 0

	for range 3 {
		increment(&a2)
	}

	a2 = 0

	for i := range 3 {
		increment(&a2)
		fmt.Print(i)
	}

	fmt.Println()
	fmt.Println(a2)

	// ex 15
	a3 := 12

	result3 := allocateInt(a3)

	fmt.Println(*result3)

	outsideLocal := localInt(a3)

	fmt.Println(outsideLocal)
}
