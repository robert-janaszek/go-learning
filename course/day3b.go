package course

import "fmt"

type address struct {
	City    string
	ZipCode string
}

type user struct {
	Name string
	Addr address
}

type embedUser struct {
	Name string
	address
}

func (a address) FullAddress() string {
	return a.City
}

func (u embedUser) FullAddress() string {
	return u.Name + " from " + u.address.FullAddress()
}

type engine struct {
	HorsePower int
}
type car struct {
	*engine
}

func (e *engine) Power() int {
	return e.HorsePower
}

type logger struct {
}

func (l logger) Log(msg string) {
	fmt.Println(msg)
}

type database struct{}

func (d database) Connect() {
	fmt.Println("Connected")
}

type server struct {
	logger
	database
}

func Day3b() {
	// ex 11
	u := user{
		Name: "Tom",
		Addr: address{
			City:    "Warsaw",
			ZipCode: "01-001",
		},
	}

	fmt.Printf("%v\n", u)
	fmt.Println(u.Addr.City)

	// ex 12

	user2 := embedUser{
		Name: "Mark",
		address: address{
			City:    "Poznan",
			ZipCode: "02-020",
		},
	}

	fmt.Println(user2.address)
	fmt.Println(user2.City)
	fmt.Println(user2.address.City)

	// ex 13

	fmt.Println("// 13")
	fmt.Println(user2.FullAddress())
	fmt.Println(user2.address.FullAddress())

	// ex 14
	c := car{}
	e := engine{
		HorsePower: 100,
	}
	c1 := car{
		engine: &e,
	}
	c2 := car{&e}
	fmt.Println(c)
	fmt.Println(c1)
	fmt.Println(c2)

	// c.Power() -- invalid memory address or nil pointer dereference

	// ex 15
	s := server{}
	s.Log("Log")
	s.Connect()
}
