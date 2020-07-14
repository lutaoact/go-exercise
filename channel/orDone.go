package main

import (
	"fmt"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-genOrDone(
		sig(3*time.Second),
		sig(5*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
		sig(2*time.Second),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}

func genOrDone(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select { // 如果从这一步就开始递归，每次只会切掉2个元素，不如从下一步开始递归，每次切掉3个，产生更少的递归
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-genOrDone(append(channels[3:], orDone)...): // 每次必须切掉2个或以上的元素
			}
		}
	}()
	return orDone
}
