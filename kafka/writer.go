package main

import (
	"context"
	"fmt"

	kafka "github.com/segmentio/kafka-go"
)

var brokers = []string{"10.1.8.95:9092", "10.1.15.25:9092", "10.1.3.161:9092"}

func main() {
	// make a writer that produces to topic-A, using the least-bytes distribution
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    "llspay__notice_telis",
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}
