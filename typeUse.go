package main

import "fmt"

type myType string

func main() {
	testCast()
}

func testCast() {
	cast(nil)
	var a string
	cast(&a)
}

func cast(i interface{}) {
	if i == nil {
		fmt.Println("nil")
		return
	}

	switch i.(type) {
	case *string:
		fmt.Println("*string")
	default:
		fmt.Println("other")
	}
}
