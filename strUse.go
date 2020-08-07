package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	parseEmptyString()
}

func parse() {
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

func parseEmptyString() {
	n, _ := strconv.Atoi("")
	fmt.Println(n)
}
