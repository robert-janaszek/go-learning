package course

import (
	"encoding/json"
	"fmt"

	"github.com/robert-janaszek/go-learning/bank"
)

type Product struct {
	ID           int     `json:"product_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	InternalCode string  `json:"-"`
	Discount     float64 `json:"discount,omitempty"`
}

type CartItem struct {
	Product  Product
	Quantity int
}
type Cart struct {
	Items []CartItem
}

func (c *Cart) AddItem(p Product, qty int) {
	item := CartItem{Product: p, Quantity: qty}
	c.Items = append(c.Items, item)
}
func (c Cart) Total() float64 {
	var total float64 = 0

	for _, item := range c.Items {
		total += float64(item.Quantity) * (item.Product.Price - item.Product.Discount)
	}

	return total
}

func Day2c() {
	// ex 16
	p := Product{ID: 1, Name: "hairdryer", Price: 300}
	j, err := json.Marshal(p)

	if err != nil {
		return
	}

	fmt.Println(string(j))

	// ex 17
	p1 := Product{
		ID:           2,
		Name:         "hairdryer 2",
		Price:        400,
		InternalCode: "secret",
		Discount:     50,
	}
	j1, err := json.Marshal(p1)

	if err != nil {
		return
	}

	fmt.Println(string(j1))

	// ex 18
	jsonData := []byte(`{"name":"Laptop", "price": 2500}`)
	var p2 Product
	err = json.Unmarshal(jsonData, &p2)

	if err != nil {
		return
	}

	fmt.Printf("%+v\n", p2)

	// ex 19
	cart := Cart{}
	cart.AddItem(p1, 1)
	cart.AddItem(p2, 3)

	fmt.Println(cart.Total())

	// ex 20

	account := bank.Account{}
	err = account.Deposit(150)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(account.Balance())

	err = account.Withdraw(100)
	if err != nil {
		fmt.Println(err)
	}

	err = account.Withdraw(100)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(account.Balance())

	// account.balance = 100 -- not allowed
}
