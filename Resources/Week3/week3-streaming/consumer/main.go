package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	// 1. Create a Kafka Reader (Consumer)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		GroupID:  "analytics-group", // THIS IS CRITICAL!
		Topic:    "order-events",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()

	log.Println("Analytics Worker started. Waiting for messages...")

	// 2. Read messages in an infinite loop
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("failed to read message: %v", err)
		}
		fmt.Printf("Worker received on Partition %d: %s = %s\n", m.Partition, string(m.Key), string(m.Value))
	}
}
