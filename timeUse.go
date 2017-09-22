package main

import (
	"fmt"
	"reflect"
	"time"
)

type Foo struct {
	Name    string
	Ports   []int
	Enabled bool
}

func main() {
	t := time.Now()
	fmt.Println(reflect.TypeOf(t))

	latency := time.Since(t)
	fmt.Println(reflect.TypeOf(latency))
	fmt.Println(latency)

	foo := Foo{Name: "gopher", Ports: []int{80, 443}, Enabled: true}
}
