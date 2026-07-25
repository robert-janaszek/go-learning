package jsonfixer

import "fmt"

func PartialParse(input string) (string, error) {
	heap := Stack{}
	var err error

	// for i := 0; i < len(input); i++ {
	// 	err = heap.AddSpecial(input[i])

	// 	if err != nil {
	// 		return "", err
	// 	}
	// }

	// fmt.Printf("%s\n", heap.heap)

	// return input, nil

	lexer := Lexer{}
	lexer.Start(input)

	nextLexem := lexer.Next()

	for nextLexem != EOF {
		err = heap.AddSpecial(nextLexem)

		if err != nil {
			return "", err
		}

		nextLexem = lexer.Next()
	}
	fmt.Printf("%s\n", heap.stack)

	return input, nil
}
