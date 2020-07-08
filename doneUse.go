package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	//readFromChan()
	sendToChan()
}

func readFromChan() {
	doWork := func(done <-chan struct{}, strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exit")
			defer close(completed)
			for {
				select {
				case v := <-strings:
					fmt.Println(v)
				case <-done:
					return
				}
			}
		}()

		return completed
	}

	done := make(chan struct{})
	completed := doWork(done, nil)

	go func() {
		time.Sleep(time.Second)
		close(done)
	}()

	<-completed
	fmt.Println("done")
}

func sendToChan() {
	newRandStream := func(done <-chan struct{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer close(randStream)
			defer fmt.Println("sendToChan exit")
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()

		return randStream
	}

	done := make(chan struct{})
	randStream := newRandStream(done)
	for i := 0; i < 3; i++ {
		fmt.Println(<-randStream)
	}
	close(done)
	time.Sleep(time.Second)
}
