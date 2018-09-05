package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/websocket"
)

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for a partner...")
	select {
	case partner <- c:
		// now handled by the other goroutine
	case p := <-partner:
		chat(p, c)
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	go io.Copy(a, b)
	io.Copy(b, a)
}

type socket struct {
	io.ReadWriter
	done chan bool
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

func socketHandler(ws *websocket.Conn) {
	s := socket{ws, make(chan bool)}
	go match(s)
	<-s.done
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func netMain() {
	listenAddr := "localhost:4000"

	go netListen()
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func netListen() {
	l, err := net.Listen("tcp", "localhost:4001")
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go match(c)
	}
}

func main() {
	netMain()
}
