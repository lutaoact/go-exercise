package main

import "fmt"

func f(left chan<- int, right <-chan int) {
	r := <-right
	fmt.Println(r)
	left <- 1 + r
}

func main() {
	const n = 10000
	leftmost := make(chan int)
	right := leftmost
	left := leftmost
	for i := 0; i < n; i++ {
		right = make(chan int)
		go f(left, right)
		left = right
	}
	go func(c chan<- int) { c <- 1 }(right)
	fmt.Println(<-leftmost)
}

/*
循环中的go f(left, right)会全部hang住，直到执行最后的 go func(c chan<- int) { c <- 1 }(right)
这个时候的right就是循环中最后定义的right = make(chan int)
收到1之后，hang住的最后一个goroutine开始执行，left中放入2，而这个left就是上一次循环中的right
也就是上一次循环中的right拿到了2，于是新的一次循环，left中放入3

边界条件是这样的：第一个goroutine退出的时候，left中放入的是2，那经过10000次之后，left中放入的就是10001
所以，最后的输出是10001
*/
