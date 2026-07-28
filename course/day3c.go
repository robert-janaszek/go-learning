package course

import "fmt"

func describe(i any) {
	fmt.Println(i)
}

func processInput(v any) {
	switch t := v.(type) {
	case int:
		fmt.Println("int: ", t)
	case string:
		fmt.Println("string: ", t)
	case bool:
		fmt.Println("bool: ", t)
	case player:
		fmt.Println("player: ", t)
	default:
		fmt.Println("unknown type")
	}
}

// Stwórz interfejs Payer.
// Stwórz strukturę CreditCard posiadającą unikalne pole CardNumber string.
// Przypisz CreditCard do zmiennej typu Payer.
// Użyj asercji typu, aby wydobyć CardNumber.

type payer interface {
	Pay(amount float64) error
}
type creditCard struct {
	CardNumber string
}

func (c *creditCard) Pay(amount float64) error {
	return nil
}

func Day3c() {
	// ex 11
	d := document{
		Name: "test.txt",
	}
	describe(d)
	describe("test")
	describe(42)

	// ex 12
	var val any = "hello Go"
	s := val.(string)
	fmt.Println(len(s))
	fmt.Println(s)

	// ex 13
	n, ok := val.(int)
	if !ok {
		fmt.Println("incorrect casting")
	}
	fmt.Println(n)

	// ex 14
	processInput(2)
	processInput("test")
	processInput(false)
	processInput(player{})

	// ex 15
	var p payer = &creditCard{
		CardNumber: "123 321",
	}
	c, ok := p.(*creditCard)

	if !ok {
		fmt.Println("cannot convert to credit card")
		return
	}
	fmt.Println("Credit card number: ", c.CardNumber)

}
