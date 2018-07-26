package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

func main() {
	//fmt.Println(MakeKey())
	testGenReqId()
}

func testGenReqId() {
	mymap := map[string]int{}
	keys := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < 1000000; i++ {
			keys <- genReqId()
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 1000000; i++ {
			keys <- genReqId()
		}
		wg.Done()
	}()

	go func() {
		for key := range keys {
			if _, ok := mymap[key]; ok {
				fmt.Printf("key = %+v dup\n", key)
				os.Exit(3)
			}
			fmt.Printf("key = %+v\n", key)
			mymap[key] = 1
		}
	}()

	wg.Wait()
	close(keys)
}

func MakeKey() string {
	var b [30]byte
	io.ReadFull(rand.Reader, b[:])
	return base64.URLEncoding.EncodeToString(b[:])
}

func genReqId() string {
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], uint32(time.Now().UnixNano()%4294967291))
	io.ReadFull(rand.Reader, b[4:])
	return base64.URLEncoding.EncodeToString(b[:])
}
