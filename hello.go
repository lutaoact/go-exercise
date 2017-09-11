package main

import "fmt"

func main() {
	fmt.Println(1<<31 - 1) //移位运算的优先级高于减法，这跟js中是不一样的
}
