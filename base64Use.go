package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func main() {
	fmt.Println(MakeKey())
}

func MakeKey() string {
	var b [30]byte
	io.ReadFull(rand.Reader, b[:])
	return base64.URLEncoding.EncodeToString(b[:])
}
