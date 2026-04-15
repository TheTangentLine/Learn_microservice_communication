package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/segmentio/kafka-go"
)

// This map acts as our lightning-fast "Read Database" (like Redis or Elasticsearch)
var readDatabase = make(map[string]string)

type PriceEvent struct {
	Item     string `json:"item"`
	NewPrice string `json:"new_price"`
}

// 1. The HTTP GET Endpoint (Lightning fast, no complex SQL joins)
func getPriceHandler(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price, exists := readDatabase[item]

	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	w.Write([]byte(fmt.Sprintf("Current price of %s is $%s\n", item, price)))
}

// 2. The Background Event Consumer
func buildReadModel() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "query-service-group",
		Topic:   "price-events",
	})
	defer reader.Close()

	log.Println("Listening to Kafka to build Read Model...")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err == nil {
			var event PriceEvent
			json.Unmarshal(m.Value, &event)

			// Update our Query Database based on the event!
			readDatabase[event.Item] = event.NewPrice
			log.Printf("Read Model Updated: %s is now $%s", event.Item, event.NewPrice)
		}
	}
}

func main() {
	// Start the Kafka consumer in a background Goroutine
	go buildReadModel()

	http.HandleFunc("/get-price", getPriceHandler)
	log.Println("Query API running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
