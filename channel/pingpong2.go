package main

import (
	"fmt"
	"time"
)

type Ball struct{ hits int }

func main() {
	done := make(chan struct{})
	table := make(chan *Ball)
	go player(done, "ping", table)
	go player(done, "pong", table)

	table <- new(Ball) // game on; toss the ball
	time.Sleep(2 * time.Second)
	<-table // game over; grab the ball
	close(done)
}

func player(done <-chan struct{}, name string, table chan *Ball) {
	for {
		select {
		case <-done:
			return
		case ball := <-table:
			ball.hits++
			fmt.Println(name, ball.hits)
			time.Sleep(100 * time.Millisecond)
			table <- ball
		}
	}
}
