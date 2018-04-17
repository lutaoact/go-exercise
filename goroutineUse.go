package main

import "fmt"

func main() {
	go func() {
		go func() {
			fmt.Println("hello 2")
		}()

		fmt.Println("hello 1")
	}()

	for {
	}
}
