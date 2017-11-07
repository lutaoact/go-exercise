package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/lutaoact/go-exercise/msgqueue"
)

func main() {
	//	ticker := time.NewTicker(100 * time.Millisecond)
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		msgqueue.Process()
	}()

	count := 1
	for t := range ticker.C {
		fmt.Printf("t = %+v\n", t)
		go msgqueue.Add(&msgqueue.Msg{
			Namespace: "library",
			RepoName:  "redis",
			Tag:       strconv.Itoa(count),
		})
		count++
	}

	for {
	}
}
