package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	//fmt.Println(MakeKey())
	testGenReqId()
}

func testGenReqId() {
	mymap := map[string]int{}
	keys := make(chan string)
	go func() {
		for i := 0; i < 1000000; i++ {
			keys <- genReqId()
		}
	}()

	go func() {
		for i := 0; i < 1000000; i++ {
			keys <- genReqId()
		}
	}()

	for {
		select {
		case key := <-keys:
			if _, ok := mymap[key]; ok {
				fmt.Printf("key = %+v dup\n", key)
				os.Exit(3)
			}
			fmt.Printf("key = %+v\n", key)
			mymap[key] = 1
		}
	}
}

func MakeKey() string {
	var b [30]byte
	io.ReadFull(rand.Reader, b[:])
	return base64.URLEncoding.EncodeToString(b[:])
}

var pid = uint32(time.Now().UnixNano() % 4294967291)

func genReqId() string {
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], pid)
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}
