package main

import "fmt"

func main() {
	fmt.Println(computeSizeTuple(131072)) // 131072 3 8
	fmt.Println(computeSizeTuple(65536))  // 65536 2 7
	fmt.Println(computeSizeTuple(65534))  // 65536 2 7
}

type sizeTuple struct {
	Size       int64
	BlockCount int64
	Group      int
}

func computeSizeTuple(size int64) *sizeTuple {
	n := size >> 10 // KB为单位
	i := 0
	for n > 0 && i < 10 {
		i += 1
		n = n >> 1
	}
	fmt.Println(i, n) // i = 8, n = 0

	ret := &sizeTuple{size, 0, 1 << (10 + i)}
	if i <= 6 {
		ret.BlockCount = 1
	} else {
		ret.BlockCount = (size >> 16) + 1
	}

	return ret
}
