package msgqueue

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Msg struct {
	Namespace string
	RepoName  string
	Tag       string
}

var queue chan *Msg = make(chan *Msg, 10)
var lock sync.Mutex

func Add(msg *Msg) {
	queue <- msg
}

func Process() {
	fmt.Println("now in process...")
	for {
		select {
		case msg := <-queue:
			log.Printf("Acquire Msg: %s %s %s", msg.Namespace, msg.RepoName, msg.Tag)
			processMsg(msg)
		default:
			log.Println("no msg in queue")
			time.Sleep(3 * time.Second)
		}
	}
}

func processMsg(msg *Msg) {
	fmt.Printf("processMsg = %+v\n", msg)
}
