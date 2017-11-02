package main

import (
	"fmt"
	"regexp"
)

func main() {
	//	searchString := "[]#:ab c"
	//	re := regexp.MustCompile(`[\-\[\]{}()*+?.,\\^$|#\s]`)
	//	searchString = re.ReplaceAllString(searchString, "\\$0")
	//	fmt.Println(searchString)

	re := regexp.MustCompile(`^[a-z0-9]+(?:[._-]?[a-z0-9]+)+$`)
	match := re.MatchString("a.b-c")
	fmt.Printf("match = %+v\n", match)
}
