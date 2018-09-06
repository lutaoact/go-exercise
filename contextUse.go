package main

import (
	"fmt"
	"time"

	"github.com/lutaoact/go-exercise/context"
)

func main() {
	contextMain()
}

func contextMain() {
	c1 := context.WithValue(context.TODO(), "hello", 1)
	c2 := context.WithValue(c1, "hello", 2)
	fmt.Println(c1.Value("hello"))
	fmt.Println(c2.Value("hello"))
	c3, cancel := context.WithTimeout(c2, 3*time.Second)
	defer cancel()
	useContext(c3)
}

func useContext(c context.Context) {
	select {
	case <-c.Done():
		fmt.Println("cancel here")
	case <-time.After(2 * time.Second):
		fmt.Println("timeout here")
	}
}
