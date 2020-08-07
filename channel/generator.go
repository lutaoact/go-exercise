package main

import (
	"fmt"
	"math/rand"
)

func main() {
	done := make(chan interface{})

	fmt.Println("out")
	values := []interface{}{1, 2, 3, 4}
	out := take(done, repeat(done, values...), 10)
	for v := range out {
		fmt.Println(v)
	}

	fmt.Println("out2")
	randFn := func() interface{} {
		return rand.Int()
	}
	out2 := take(done, repeatFn(done, randFn), 10)
	for v := range out2 {
		fmt.Println(v)
	}
}

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			for _, v := range values {
				select {
				case out <- v:
				case <-done:
					return
				}
			}
		}
	}()
	return out
}

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case out <- fn():
			case <-done:
				return
			}
		}
	}()
	return out
}

func take(done <-chan interface{}, in <-chan interface{}, num int) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for i := 0; i < num; i++ {
			select {
			case out <- <-in:
			case <-done:
				return
			}
		}
	}()
	return out
}
