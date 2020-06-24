package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	randN()
}

func randN() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		index := rand.Int63n(10)
		fmt.Printf("index = %+v\n", index)
	}
}
