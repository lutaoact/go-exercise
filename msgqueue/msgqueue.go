package msgqueue

import (
	"fmt"
	"log"
	"time"
)

type Msg struct {
	Namespace string
	RepoName  string
	Tag       string
}

var queue chan *Msg = make(chan *Msg, 512)

func Add(msg *Msg) {
	queue <- msg
}

func Process(handler func(*Msg)) {
	fmt.Println("now in process...")
	for {
		select {
		case msg := <-queue:
			log.Printf("Acquire Msg: %s %s %s", msg.Namespace, msg.RepoName, msg.Tag)
			handler(msg)
		default:
			log.Println("no msg in queue")
			time.Sleep(3 * time.Second)
		}
	}
}
