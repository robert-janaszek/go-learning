package course

import (
	"fmt"
)

func send(ch chan<- string) {
	ch <- "ping"
}

// Zadanie 10: Bezpieczne sprawdzanie czy kanał jest otwarty
// Odbierz wartość z zamkniętego kanału używając składni dwuargumentowej:
// val, ok := <-ch. Sprawdź, jaką wartość ma val oraz ok dla
// otwartego i zamkniętego kanału.

func producer(ch chan<- int) {
	for i := 1; i < 6; i++ {
		ch <- i
	}
	close(ch)
}

func Day6b() {
	// ex 6
	ch := make(chan string)
	go send(ch)
	msg := <-ch
	fmt.Println(msg)
	// ex 7
	// ch <- "data" // all goroutines are asleep - deadlock!

	// ex 8
	chb := make(chan int, 2)
	chb <- 1
	chb <- 2
	// chb <- 3 // deadlock

	// ex 9
	chc := make(chan int)

	go producer(chc)

	val, ok := <-chc
	fmt.Println(val, ok)
	for val = range chc {
		fmt.Println(val)
	}

	// ex 10
	val, ok = <-chc
	fmt.Println(val, ok)

}
