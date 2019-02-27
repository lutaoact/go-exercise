package main

import "fmt"

type data struct {
	name string
}

func main() {
	m := map[string]*data{
		"x": {"Tom"},
	}

	m["x"].name = "Jerry" // 直接修改 m["x"] 中的字段
	fmt.Println(m["x"])   // &{Jerry}
}
