package days

import (
	"fmt"
	"sync"
)

func sayHello(wg *sync.WaitGroup, name string) {
	defer wg.Done()
	fmt.Printf("Hi %s!\n", name)
}

func printI(wg *sync.WaitGroup, i int) {
	defer wg.Done()
	fmt.Println("i: ", i)
}

func inc(wg *sync.WaitGroup, mu *sync.Mutex, i *int) {
	defer wg.Done()
	defer mu.Unlock()

	mu.Lock()
	*i++
}

func Day7a() {
	// ex 2
	var wg sync.WaitGroup
	wg.Add(1)
	// ex 1
	go sayHello(&wg, "Robert")
	wg.Wait()

	// ex 3
	wg.Add(10)
	for i := 0; i < 5; i++ {
		go func() { defer wg.Done(); fmt.Println(i) }()
		go printI(&wg, i)
	}
	wg.Wait()

	// ex 4
	counter := 0
	var mu sync.Mutex
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go inc(&wg, &mu, &counter)
	}
	wg.Wait()
	fmt.Print(counter) // 96, 100 with mutex
}
