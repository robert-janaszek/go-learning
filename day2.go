package main

import (
	"fmt"
)

type Book struct {
	Title  string
	Author string
	Pages  int
	IsRead bool
}

func (b Book) Summary() string {
	// return "\"" + b.Title + "\" by " + b.Author + " (" + strconv.Itoa(b.Pages) + " pages)"
	return fmt.Sprintf("\"%s\" by %s (%d pages)", b.Title, b.Author, b.Pages)
}

func (b *Book) MarkAsRead() {
	b.IsRead = true
}

type BookWithTags struct {
	Title  string
	Author string
	Pages  int
	IsRead bool
	Tags   []string
}

func NewBook(title string, author string, pages int) *Book {
	return &Book{
		Title:  title,
		Author: author,
		Pages:  pages,
	}
}

type Celsius float64

func (c Celsius) ToFahrenheit() float64 {
	return float64(c)*9/5 + 32
}

func (c *Celsius) Add(value float64) {
	// *c = Celsius(float64(*c) + value) -- ok
	*c += Celsius(value)
}

func day2() {
	// ex 1
	book := Book{
		Title:  "De revolutionibus orbium coelestium",
		Author: "Nicolas Copernicus",
		Pages:  392,
		IsRead: true,
	}
	fmt.Printf("%+v\n", book)

	// ex 2
	book2 := Book{
		"Philosophiae Naturalis Principia Mathematica",
		"Isaac Newton",
		423,
		false,
	}
	fmt.Printf("%+v\n", book2)

	book3 := Book{}
	book3.Title = "Special theory of relativity"
	book3.Author = "Albert Einstein"
	book3.Pages = 543
	book3.IsRead = true

	fmt.Printf("%+v\n", book3)

	// ex 3

	book4a := NewBook("On the Origin of Species", "Charles Darwin", 500)
	book4b := *NewBook("On the Origin of Species", "Charles Darwin", 500) // copies value - mistake

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

	book5a := Book{
		Title:  "De revolutionibus orbium coelestium",
		Author: "Nicolas Copernicus",
		Pages:  392,
		IsRead: true,
	}
	book5b := Book{
		Title:  "De revolutionibus orbium coelestium",
		Author: "Nicolas Copernicus",
		Pages:  392,
		IsRead: true,
	}

	fmt.Println(book5a == book5b)

	// ex 5

	book6a := BookWithTags{}
	book6b := BookWithTags{}

	fmt.Println(book6a, book6b)

	// fmt.Println(book6a == book6b) // Do not compile

	// ex 6
	fmt.Println(book.Summary())

	// ex 7
	book2.MarkAsRead()
	fmt.Printf("%+v\n", book2)

	// ex 8
	book7 := Book{Title: "Go in Action"}
	book7.MarkAsRead()
	// (&book7).MarkAsRead() - not needed
	fmt.Printf("%+v\n", book7)

	// ex 9
	temp := Celsius(100)
	fmt.Println(temp.ToFahrenheit())

	// ex 10
	temp.Add(10)
	fmt.Println(temp)
}
