package main

import (
	"fmt"
	"os"
)

func main() {
	fi, err := os.Stat("/tmp")
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("fi = %+v\n", fi)
	fmt.Printf("fi.IsDir() = %+v\n", fi.IsDir())

	fi, err = os.Stat("/tmp/tmp.txt")
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("fi = %+v\n", fi)
	fmt.Printf("fi.IsDir() = %+v\n", fi.IsDir())
}
