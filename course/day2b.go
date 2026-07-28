package course

import "fmt"

type Address struct {
	City    string
	ZipCode string
}

type User struct {
	Name string
	Addr Address
}

type EmbedUser struct {
	Name string
	Address
}

func (a Address) FullAddress() string {
	return a.City
}

func (u EmbedUser) FullAddress() string {
	return u.Name + " from " + u.Address.FullAddress()
}

type Engine struct {
	HorsePower int
}
type Car struct {
	*Engine
}

func (e *Engine) Power() int {
	return e.HorsePower
}

type Logger struct {
}

func (l Logger) Log(msg string) {
	fmt.Println(msg)
}

type Database struct{}

func (d Database) Connect() {
	fmt.Println("Connected")
}

type Server struct {
	Logger
	Database
}

func Day2b() {
	// ex 11
	user := User{
		Name: "Tom",
		Addr: Address{
			City:    "Warsaw",
			ZipCode: "01-001",
		},
	}

	fmt.Printf("%v\n", user)
	fmt.Println(user.Addr.City)

	// ex 12

	user2 := EmbedUser{
		Name: "Mark",
		Address: Address{
			City:    "Poznan",
			ZipCode: "02-020",
		},
	}

	fmt.Println(user2.Address)
	fmt.Println(user2.City)
	fmt.Println(user2.Address.City)

	// ex 13

	fmt.Println("// 13")
	fmt.Println(user2.FullAddress())
	fmt.Println(user2.Address.FullAddress())

	// ex 14
	c := Car{}
	e := Engine{
		HorsePower: 100,
	}
	c1 := Car{
		Engine: &e,
	}
	c2 := Car{&e}
	fmt.Println(c)
	fmt.Println(c1)
	fmt.Println(c2)

	// c.Power() -- invalid memory address or nil pointer dereference

	// ex 15
	server := Server{}
	server.Log("Log")
	server.Connect()
}
