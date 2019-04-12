package main

import (
	"fmt"

	"github.com/buger/jsonparser"
)

var data = []byte(`{"hello":{"good":"ok","only":true}}`)

func main() {
	//fmt.Println(jsonparser.Get(data, "hello", "only"))
	//fmt.Println(string(data))
	//data = jsonparser.Delete(data, "hello", "only") //delete之后，必须赋值给变量，否则会出奇怪的问题
	//fmt.Println(string(data))
	//fmt.Println(string(jsonparser.Delete(data, "hello", "good")))

	fmt.Println(string(data))
	fmt.Println(string(jsonparser.Delete(data, "hello", "ok")))
	//	fmt.Println(string(jsonparser.Delete(data, "hello", "only")))
	//	fmt.Println(string(jsonparser.Delete(data, "hello", "good")))
}
