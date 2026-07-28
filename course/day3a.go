package course

import "fmt"

type stringer interface {
	String() string
}

type namedUser struct {
	Name string
	Age  int
}

func (u namedUser) String() string {
	return fmt.Sprintf("%s, %d", u.Name, u.Age)
}

var _ stringer = namedUser{}

func Day3a() {
	// ex 1
	mark := namedUser{
		Name: "Mark",
		Age:  32,
	}

	fmt.Println(mark.String())
}
