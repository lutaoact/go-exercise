package main

import "fmt"

// 正确示例
func main() {
	defer func() {
		fmt.Println("recovered:", recover())
	}()
	panic("not good")
	println("hhh")
}

func main1() {
	recover()         // 什么都不会捕捉
	panic("not good") // 发生 panic，主程序退出
	recover()         // 不会被执行
	println("ok")
}

func main2() {
	defer func() {
		doRecover()
	}()
	panic("not good")
}

func doRecover() {
	fmt.Println("recovered:", recover())
}
