package course

import (
	"fmt"
	"sync"
	"time"
)

func worker1(wg *sync.WaitGroup, jobs <-chan int, results chan<- int) {
	defer wg.Done()
	for job := range jobs {
		results <- job * 2
		time.Sleep(20 * time.Millisecond)
	}
}

func ToSlice(c chan int) []int {
	s := make([]int, 0)
	for i := range c {
		s = append(s, i)
	}
	return s
}

func merge(ch ...<-chan int) chan int {
	var result chan int = make(chan int, 1000)
	for _, c := range ch {
		for v := range c {
			result <- v
		}
	}

	close(result)

	return result
}

func fetch1(ch chan<- string) {
	time.Sleep(190 * time.Millisecond)
	ch <- "fetch 1"
}
func fetch2(ch chan<- string) {
	time.Sleep(200 * time.Millisecond)
	ch <- "fetch 2"
}
func fetch3(ch chan<- string) {
	time.Sleep(210 * time.Millisecond)
	ch <- "fetch 3"
}

func ddos(ch chan string) {
	for i := 0; i < 10; i++ {
		ch <- "DDoS"
	}
}

func process(ch chan struct{}) {
	defer close(ch)
	time.Sleep(1200 * time.Millisecond)
	// ch <- struct{}{}
}

func Day7d() {
	// ex 16
	wg := sync.WaitGroup{}
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	wg.Add(3)
	jobs <- 1
	jobs <- 2
	jobs <- 3
	jobs <- 4
	jobs <- 5
	jobs <- 6
	jobs <- 7
	jobs <- 8
	jobs <- 9
	jobs <- 10

	go worker1(&wg, jobs, results)
	go worker1(&wg, jobs, results)
	go worker1(&wg, jobs, results)

	close(jobs)
	wg.Wait()
	close(results)

	fmt.Println(ToSlice(results))

	// ex 17
	jobs1 := make(chan int, 100)
	jobs2 := make(chan int, 100)
	jobs3 := make(chan int, 100)
	results1 := make(chan int, 100)
	results2 := make(chan int, 100)
	results3 := make(chan int, 100)

	wg.Add(3)
	go worker1(&wg, jobs1, results1)
	go worker1(&wg, jobs2, results2)
	go worker1(&wg, jobs3, results3)

	for i := 0; i < 20; i++ {
		workerNum := i % 3

		switch workerNum {
		case 0:
			jobs1 <- i
		case 1:
			jobs2 <- i
		case 2:
			jobs3 <- i
		}
	}

	close(jobs1)
	close(jobs2)
	close(jobs3)
	wg.Wait()
	close(results1)
	close(results2)
	close(results3)

	allResults := merge(results1, results2, results3)
	fmt.Println(ToSlice(allResults))

	// ex 18
	var chanHttp = make(chan string, 3)
	go fetch1(chanHttp)
	go fetch2(chanHttp)
	go fetch3(chanHttp)

	result := <-chanHttp
	fmt.Println(result)

	// ex 19
	var throttledCh = make(chan string, 1)

	go ddos(throttledCh)

	ticker := time.Tick(200 * time.Millisecond)

	for i := 0; i < 10; i++ {
		<-ticker
		v := <-throttledCh
		fmt.Println(v)
	}

	// ex 20
	done := make(chan struct{})

	go process(done)

	<-done
}
