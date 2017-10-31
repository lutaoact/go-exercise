package main

import (
	"fmt"
	"regexp"
)

func main() {
	searchString := "[]#:ab c"
	re := regexp.MustCompile(`[\-\[\]{}()*+?.,\\^$|#\s]`)
	searchString = re.ReplaceAllString(searchString, "\\$0")
	fmt.Println(searchString)
}
