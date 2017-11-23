package main

import (
	"fmt"

	"github.com/vmihailenco/msgpack"
)

func main() {
	ExampleMarshal()
}

func ExampleMarshal() {
	type Item struct {
		Foo string
	}

	b, err := msgpack.Marshal(&Item{Foo: "bar"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("string(b) = %+v\n", string(b))

	var item Item
	err = msgpack.Unmarshal(b, &item)
	if err != nil {
		panic(err)
	}
	fmt.Println(item.Foo)
	// Output: bar
}
