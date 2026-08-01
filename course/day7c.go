package course

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func producer1(ch chan<- int, wait time.Duration, val int) {
	time.Sleep(wait)
	ch <- val
}

func produce(ch chan<- int) {
	ch <- 11
}

func consume(ch <-chan int) {
	val := <-ch
	fmt.Println(val)
}

func worker(ctx context.Context, wg *sync.WaitGroup, work <-chan int) {
	defer wg.Done()
	for {
		select {
		case val := <-work:
			fmt.Println("working ", val)
			time.Sleep(30 * time.Millisecond)

		case <-ctx.Done():
			fmt.Println("finishing work")
			return
		}
	}
}

func nonBlockingWorker(ctx context.Context, wg *sync.WaitGroup, work <-chan int) {
	defer wg.Done()
	for {
		select {
		case val := <-work:
			fmt.Println("working ", val)
			time.Sleep(30 * time.Millisecond)
		case <-ctx.Done():
			fmt.Println("finishing work")
			return
		default:
			fmt.Println("idle")
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func Day7c() {
	//ex 11
	ch1 := make(chan int)
	ch2 := make(chan int)
	go producer1(ch1, 400*time.Millisecond, 1)
	go producer1(ch2, 250*time.Millisecond, 2)

	select {
	case val := <-ch1:
		fmt.Println(val)
	case val := <-ch2:
		fmt.Println(val)
	}

	// ex 12
	ch3 := make(chan int)
	go producer1(ch3, 3*time.Second, 3)

	select {
	case val := <-ch3:
		fmt.Println(val)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}

	// ex 13
	ch4 := make(chan int)
	go produce(ch4)
	consume(ch4)

	// ex 14
	ch5 := make(chan int, 5)
	ch5 <- 1
	ch5 <- 2
	ch5 <- 3
	ch5 <- 4
	ch5 <- 5

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	wg := sync.WaitGroup{}

	wg.Add(1)
	go worker(ctx, &wg, ch5)
	wg.Wait()

	// ex 15
	ch6 := make(chan int, 5)
	ch6 <- 1
	ch6 <- 2
	ch6 <- 3
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	wg.Add(1)
	go nonBlockingWorker(ctx, &wg, ch6)
	wg.Wait()
}
