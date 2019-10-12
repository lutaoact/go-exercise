package main

import (
	"fmt"
)

func main() {
	printRawBytes()
}

func IsBigEndian() bool {
	var i int32 = 0x12345678
	var b byte = byte(i)
	if b == 0x12 {
		return true
	}

	return false
}

func mainIsBigEndian() {
	if IsBigEndian() {
		fmt.Println("大端序")
	} else {
		fmt.Println("小端序")
	}
}

func printRawBytes() {
	bs := []byte{0x00, 0xfd, 0x12}
	for _, n := range bs {
		fmt.Printf("%08b\n", n) // prints 00000000 11111101
	}
}
