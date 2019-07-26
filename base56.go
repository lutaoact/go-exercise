package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

const (
	// 在base64的基础上，减少了+/两个符号以及肉眼难以分辨 1 l I o O 0，即base56
	alphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"

	alphabetIdx0 = '1'
)

var bigRadix = big.NewInt(56)
var bigZero = big.NewInt(0)

func Encode(b []byte) string {
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0, len(b)*2)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		answer = append(answer, alphabet[mod.Int64()])
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, alphabetIdx0)
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
}

func main() {
	var b [7]byte
	io.ReadFull(rand.Reader, b[:])
	result := Encode(b[:])
	fmt.Println(result)
}
