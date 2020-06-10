package main

import "fmt"

type retForChan struct {
	Ret interface{}
	Err error
}

type MyFunc func(args ...interface{}) (interface{}, error)

func asyncCall(out chan<- *retForChan, fn MyFunc, params ...interface{}) {
	go func() {
		ret := new(retForChan)
		ret.Ret, ret.Err = fn(params...)
		out <- ret
	}()
}

func firstFunc(args ...interface{}) (interface{}, error) {
	return args[0], nil
}

func main() {
	out := make(chan *retForChan)
	asyncCall(out, firstFunc, 5)
	fmt.Println(<-out)
}
