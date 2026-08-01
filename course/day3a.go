package course

import (
	"fmt"
)

type book struct {
	Title  string
	Author string
	Pages  int
	IsRead bool
}

func (b book) Summary() string {
	// return "\"" + b.Title + "\" by " + b.Author + " (" + strconv.Itoa(b.Pages) + " pages)"
	return fmt.Sprintf("\"%s\" by %s (%d pages)", b.Title, b.Author, b.Pages)
}

func (b *book) MarkAsRead() {
	b.IsRead = true
}

type bookWithTags struct {
	Title  string
	Author string
	Pages  int
	IsRead bool
	Tags   []string
}

func newBook(title string, author string, pages int) *book {
	return &book{
		Title:  title,
		Author: author,
		Pages:  pages,
	}
}

type celsius float64

func (c celsius) ToFahrenheit() float64 {
	return float64(c)*9/5 + 32
}

func (c *celsius) Add(value float64) {
	// *c = celsius(float64(*c) + value) -- ok
	*c += celsius(value)
}

func Day3a() {
	// ex 1
	b := book{
		Title:  "De revolutionibus orbium coelestium",
		Author: "Nicolas Copernicus",
		Pages:  392,
		IsRead: true,
	}
	fmt.Printf("%+v\n", b)

	// ex 2
	book2 := book{
		"Philosophiae Naturalis Principia Mathematica",
		"Isaac Newton",
		423,
		false,
	}
	fmt.Printf("%+v\n", book2)

	book3 := book{}
	book3.Title = "Special theory of relativity"
	book3.Author = "Albert Einstein"
	book3.Pages = 543
	book3.IsRead = true

	fmt.Printf("%+v\n", book3)

	// ex 3

	book4a := newBook("On the Origin of Species", "Charles Darwin", 500)
	book4b := *newBook("On the Origin of Species", "Charles Darwin", 500) // copies value - mistake

	fmt.Printf("%+v\n", book4a)
	fmt.Printf("%+v\n", book4b)

	// ex 4
	config := struct {
		ConfigName string
		Port       int
	}{
		ConfigName: "test.yml",
		Port:       1234,
	}
	fmt.Printf("%+v\n", config)

	book5a := book{
		Title:  "De revolutionibus orbium coelestium",
		Author: "Nicolas Copernicus",
		Pages:  392,
		IsRead: true,
	}
	book5b := book{
		Title:  "De revolutionibus orbium coelestium",
		Author: "Nicolas Copernicus",
		Pages:  392,
		IsRead: true,
	}

	fmt.Println(book5a == book5b)

	// ex 5

	book6a := bookWithTags{}
	book6b := bookWithTags{}

	fmt.Println(book6a, book6b)

	// fmt.Println(book6a == book6b) // Do not compile

	// ex 6
	fmt.Println(b.Summary())

	// ex 7
	book2.MarkAsRead()
	fmt.Printf("%+v\n", book2)

	// ex 8
	book7 := book{Title: "Go in Action"}
	book7.MarkAsRead()
	// (&book7).MarkAsRead() - not needed
	fmt.Printf("%+v\n", book7)

	// ex 9
	temp := celsius(100)
	fmt.Println(temp.ToFahrenheit())

	// ex 10
	temp.Add(10)
	fmt.Println(temp)
}
