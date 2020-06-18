package main

import (
	"bytes"
	"fmt"
)

func main() {
	cmp()
	//printRawBytes()
}

func cmp() {
	a1 := []byte{0x20, 0x30}
	a2 := []byte{0x20, 0x30, 0x25}
	fmt.Println(bytes.Compare(a1, a2))
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
	fmt.Printf("bs = %08b\n", bs)
}
