package main

import "fmt"

type Book struct {
	Title  string
	Author string
	Pages  int
	IsRead bool
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

	fmt.Print(book6a, book6b)

	// fmt.Println(book6a == book6b) // Do not compile
}
