package main

import (
	"log"
	"os"
	"runner"
	"time"
)

const timeout = 3 * time.Second

func main() {
	log.Println("Starting work.")

	r := runner.New(timeout)
	r.Add(createTask(), createTask(), createTask())

	if err := r.Start(); err != nil {
		switch err {
		case runner.ErrTimeout:
			log.Println("Terminationg due to timeout.")
			os.Exit(1)
		case runner.ErrInterrupt:
			log.Println("Terminationg due to interrupt.")
			os.Exit(2)
		}
	}
	log.Println("Process ended.")
}

func createTask() func(int) {
	return func(id int) {
		log.Printf("Processor - Task #%d.", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}
