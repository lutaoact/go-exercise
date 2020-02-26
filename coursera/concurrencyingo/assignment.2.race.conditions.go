package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

/*
Two goroutines increase one variable count, each increase 10000, the result should be 20000. Actually, the result may be less than 20000.
When run `count++`, first load count, run count + 1, then save result to count, two goroutine may load count at the same time, then save result at the same time. One result will override the other one. This is explanation for race conditions of variable count.
*/
func main() {
	count := 0
	wg.Add(2)

	go func() { //goroutine1
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			count++
		}
	}()

	go func() { //goroutine2
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			count++
		}
	}()

	wg.Wait()
	fmt.Println(count)
}
