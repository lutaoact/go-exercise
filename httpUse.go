package main

import (
	"fmt"
	"strings"
)

type Header map[string][]string

func main() {
	header := Header{
		"Content-Type": []string{"application/json"},
	}
	for k, v := range header {
		fmt.Println(k, strings.Join(v, ""))
	}
}
