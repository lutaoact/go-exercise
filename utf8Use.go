package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	testfor2()
}

func testCount() {
	str := "Hello, 世界"
	fmt.Println("bytes =", len(str))
	fmt.Println("runes =", utf8.RuneCountInString(str))
}

func testfor() {
	data := "A\xfe\x02\xff\x04"
	for _, v := range data {
		fmt.Printf("%#x ", v)
	}
	//prints: 0x41 0xfffd 0x2 0xfffd 0x4 (not ok)

	fmt.Println()
	for _, v := range []byte(data) {
		fmt.Printf("%#x ", v)
	}
	//prints: 0x41 0xfe 0x2 0xff 0x4 (good)
}

func testfor2() {
	data := "A中国éA"
	for k, v := range data {
		fmt.Printf("%d: %#x\n", k, v)
	}
	//range得到的k并不是字符的位置，而是字符的字节序列中第一个字节的位置，结果为0 1 4 7 9，打出来的是unicode码点
	//0: 0x41
	//1: 0x4e2d
	//4: 0x56fd
	//7: 0xe9
	//9: 0x41
}
