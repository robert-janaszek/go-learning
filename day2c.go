package main

import (
	"encoding/json"
	"fmt"
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

func (c *Cart) AddItem(p Product, qty int) {}
func (c Cart) Total() float64 {
	return 0
}

func day2c() {
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
}
