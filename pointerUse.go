package main

import "fmt"

func main() {
	boolPointer()
}

func True() *bool {
	a := true
	return &a
}

func boolPointer() {
	var a *bool
	fmt.Printf("a == nil = %+v\n", a == nil)
	fmt.Printf("a = %+v\n", a)
	a = True()
	fmt.Printf("a = %+v\n", *a)
}
