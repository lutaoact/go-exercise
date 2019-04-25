package main

import (
	"flag"
	"fmt"
)

func main() {
	main1()
}

type b struct {
	c int
	d string
}

type testST struct {
	a string
	b b
}

func main1() {
	t := &testST{}
	flag.StringVar(&t.a, "a", "", "this is a")
	flag.IntVar(&t.b.c, "c", 2, "this is c")
	flag.StringVar(&t.b.d, "d", "", "this is d")
	flag.Parse()
	fmt.Printf("t = %+v\n", t)
}
