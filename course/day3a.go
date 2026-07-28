package course

import "fmt"

type stringer interface {
	String() string
}

type namedUser struct {
	Name string
	Age  int
}

type namedBook struct {
	Title  string
	Author string
}

func (u namedUser) String() string {
	return fmt.Sprintf("%s, %d", u.Name, u.Age)
}

func (b namedBook) String() string {
	return fmt.Sprintf("\"%s\" by %s", b.Title, b.Author)
}

var _ stringer = namedUser{}
var _ stringer = namedBook{}

func printInfo(s stringer) {
	fmt.Println(s.String())
}

func Day3a() {
	// ex 1
	mark := namedUser{
		Name: "Mark",
		Age:  32,
	}

	fmt.Println(mark.String())

	// ex 2
	printInfo(mark)

	// ex 3

	orbium := namedBook{
		Title:  "De revolutionibus orbium coelestium",
		Author: "Nicolas Copernicus",
	}

	printInfo(orbium)

	// ex 4

	var s stringer
	if s == nil {
		fmt.Println("s is nil")
		// s.String() -- panic: runtime error: invalid memory address or nil pointer dereference
	}

	// ex 5
	items := []stringer{mark, orbium}
	for _, item := range items {
		fmt.Println(item.String())
	}
}
