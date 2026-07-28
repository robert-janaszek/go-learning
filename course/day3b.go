package course

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type saver interface {
	Save() error
}

type document struct {
	Name string
}

func (d *document) Save() error {
	fmt.Println("Saving document " + d.Name)
	return nil
}

type note struct {
	note string
}

func (n note) Save() error {
	return nil
}

// pointer receiver implements method only in pointer mode, contrary to value receiver
var _ saver = &document{}

// ex 7
// var _ saver = document{} -- cannot use document{} (value of struct type document) as saver value in variable declaration: document does not implement saver (method Save has pointer receiver)

func writeHello(w io.Writer) {
	msg := []byte("Hello Go")
	_, err := w.Write(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func Day3b() {
	// ex 6
	doc := document{
		Name: "CV.pdf",
	}

	err := doc.Save()
	if err != nil {
		fmt.Println(err)
	}

	// ex 8
	var _ saver = note{}
	var _ saver = &note{}

	// ex 9
	writeHello(os.Stdout)
	buffer := bytes.Buffer{}
	writeHello(&buffer)

	fmt.Println(buffer.String())

	// ex 10
	// os.Stdout - done above
	file, err := os.Create("test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	writeHello(file)
}
