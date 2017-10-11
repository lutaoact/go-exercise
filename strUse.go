package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	page, err := strconv.Atoi("2")
	fmt.Println(page)
	if err != nil || page <= 0 {
		fmt.Printf("err = %+v\n", err)

		page = 10
	}
	fmt.Println(page)

	strs := strings.Split("test", "/")
	fmt.Printf("%v\n", len(strs))
	fmt.Printf("%q\n", strs)
}
