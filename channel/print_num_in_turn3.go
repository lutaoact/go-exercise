package main

import "fmt"

func main() {
	counter1 := createCounter(1)
	counter2 := createCounter(2)

	for i := 0; i < 5; i++ {
		fmt.Println(<-counter1)
		fmt.Println(<-counter2)
	}
}

func createCounter(i int) <-chan int {
	next := make(chan int)
	go func() {
		for {
			next <- i
			i += 2
		}
	}()
	return next
}
