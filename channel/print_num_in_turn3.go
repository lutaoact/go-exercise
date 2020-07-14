package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan struct{})
	counter1 := createCounter(1, done)
	counter2 := createCounter(2, done)

	for i := 0; i < 5; i++ {
		fmt.Println(<-counter1)
		fmt.Println(<-counter2)
	}
	close(done)
	time.Sleep(time.Second)
}

func createCounter(i int, done chan struct{}) <-chan int {
	next := make(chan int)
	go func() {
		for {
			select {
			case next <- i:
				i += 2
			case <-done:
				fmt.Println("here")
				return
			}
		}
	}()
	return next
}
