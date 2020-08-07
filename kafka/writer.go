package main

import (
	"context"
	"fmt"

	kafka "github.com/segmentio/kafka-go"
)

var brokers = []string{"host:port"}

func main() {
	// make a writer that produces to topic-A, using the least-bytes distribution
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    "topic_1",
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   nil,
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   nil,
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   nil,
			Value: []byte("Two!"),
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}
