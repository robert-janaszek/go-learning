package days

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(count int) *Semaphore {
	ch := make(chan struct{}, count)

	for range count {
		ch <- struct{}{}
	}

	return &Semaphore{
		ch: ch,
	}
}

func (s *Semaphore) Acquire() func() {
	<-s.ch

	return func() {
		s.ch <- struct{}{}
	}
}

func Day7e() {
	semaphore := NewSemaphore(3)
	wg := sync.WaitGroup{}

	wg.Add(4)

	go func() {
		defer wg.Done()
		release := semaphore.Acquire()
		defer release()
		fmt.Println("go 1")
		time.Sleep(3 * time.Second)
	}()
	go func() {
		defer wg.Done()
		release := semaphore.Acquire()
		defer release()
		fmt.Println("go 2")
		time.Sleep(3 * time.Second)
	}()
	go func() {
		defer wg.Done()
		release := semaphore.Acquire()
		defer release()
		fmt.Println("go 3")
		time.Sleep(3 * time.Second)
	}()
	go func() {
		defer wg.Done()
		release := semaphore.Acquire()
		defer release()
		fmt.Println("go 4")
		time.Sleep(3 * time.Second)
	}()

	wg.Wait()
}
