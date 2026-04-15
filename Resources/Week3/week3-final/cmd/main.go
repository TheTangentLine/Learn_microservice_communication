package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

func updatePriceHandler(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")

	// 1. Create the Event (We are storing the FACT that a price changed)
	eventPayload := fmt.Sprintf(`{"item": "%s", "new_price": "%s"}`, item, price)

	// 2. Publish to Kafka (Event Sourcing)
	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(item), // Keep events for the same item in order!
			Value: []byte(eventPayload),
		},
	)
	if err != nil {
		http.Error(w, "Failed to publish event", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Command Accepted: %s price update queued.\n", item)))
}

func main() {
	writer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "price-events",
		Balancer: &kafka.Hash{},
	}
	defer writer.Close()

	http.HandleFunc("/update-price", updatePriceHandler)
	log.Println("Command API running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
