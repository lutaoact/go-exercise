package main

import (
	"fmt"
	"time"
)

func main() {
	main3()
}

func main3() {
	//在值为nil的channel上发送和接收数据将永久阻塞，利用这个死锁的特性，在select中动态的打开和关闭case语句块
	inCh := make(chan int)
	outCh := make(chan int)
	go func() {
		var in <-chan int = inCh
		var out chan<- int
		var val int
		for {
			select {
			case out <- val:
				println("---------")
				out = nil
				in = inCh
			case val = <-in:
				println("+++++++++")
				in = nil
				out = outCh
			}
		}
	}()

	go func() {
		for v := range outCh {
			println("processed:", v)
		}
	}()
	inCh <- 1
	inCh <- 2
	time.Sleep(3 * time.Second)
}

func main2() {
	ch := make(chan string)

	go func() {
		for m := range ch {
			fmt.Println("Processed:", m)
			time.Sleep(1 * time.Second) // 模拟需要长时间运行的操作
			fmt.Println("done *")
		}
	}()

	ch <- "cmd.1"
	ch <- "cmd.2"               // 不会被接收处理
	ch <- "cmd.3"               // 不会被接收处理
	time.Sleep(1 * time.Second) // 模拟需要长时间运行的操作
	fmt.Println("done")
}

func main1() {
	ch := make(chan int)
	for i := 0; i < 3; i++ {
		go func(idx int) {
			ch <- idx
		}(i)
	}

	fmt.Println(<-ch)           // 输出第一个发送的值
	fmt.Println(<-ch)           // 输出第一个发送的值
	fmt.Println(<-ch)           // 输出第一个发送的值
	time.Sleep(2 * time.Second) // 模拟做其他的操作
}
