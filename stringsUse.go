package main

import (
	"fmt"
	"strings"
)

func main() {
	split()
}

func split() {
	s := "abc"
	parts := strings.SplitN(s, "=", 2)
	fmt.Printf("parts = %+v\n", parts)
}
