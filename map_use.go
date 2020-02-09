package main

import "fmt"

type data struct {
	name string
}

func main() {
	getFromNilMap()
}

func main1() {
	m := map[string]*data{
		"x": {"Tom"},
	}

	m["x"].name = "Jerry" // 直接修改 m["x"] 中的字段
	fmt.Println(m["x"])   // &{Jerry}
}

func notExistKey() {
	a := make(map[string]string)
	fmt.Println("1" + a["xxx"] + "2")
}

func getFromNilMap() {
	var m map[string]int
	fmt.Println(m == nil) // true
	fmt.Println(m["a"])   // 0
	m = make(map[string]int)
	fmt.Println(m == nil) // false
	fmt.Println(m["a"])   // 0

	var n []int
	fmt.Println(n == nil) // true
	//fmt.Println(n[0]) // panic: runtime error: index out of range [0] with length 0
}
