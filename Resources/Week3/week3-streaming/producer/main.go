package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// 1. Create a Kafka Writer (Producer)
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "order-events",
		Balancer: &kafka.Hash{}, // Routes messages with same key to same partition
	}
	defer w.Close()

	// 2. Publish 10 messages in a loop
	for i := 1; i <= 10; i++ {
		// We use the user_id as the Key.
		// All messages for User_99 will ALWAYS go to the same partition!
		key := []byte("user_99")
		value := []byte(fmt.Sprintf("Order %d placed by User 99", i))

		err := w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   key,
				Value: value,
			},
		)
		if err != nil {
			log.Fatalf("failed to write messages: %v", err)
		}

		log.Printf("Sent: %s", string(value))
		time.Sleep(1 * time.Second)
	}
}
