package main

import "fmt"

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func sqMain() {
	c := gen(2, 3)
	out := sq(c)
	fmt.Println(<-out)
	fmt.Println(<-out)
}

func main() {
	sqMain()
}
