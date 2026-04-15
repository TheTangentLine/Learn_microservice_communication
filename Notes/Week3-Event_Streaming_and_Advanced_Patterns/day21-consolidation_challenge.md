### **Day 21: Week 3 Consolidation Project**

Today, we are building a mini CQRS and Event Sourced architecture. We are going to separate our writes from our reads!

#### **The Architecture**

1.  **The Command API (Producer):** Receives HTTP requests to create products and change prices. It does _not_ save to a database. It only publishes events to Kafka.
2.  **Kafka (The Event Store):** Holds the immutable history of all price changes.
3.  **The Read API (Consumer + DB):** Listens to Kafka, updates an in-memory map (acting as our NoSQL Query Database), and serves blazing-fast HTTP GET requests to users.

#### **1. Project Setup**

Create a new folder: `week3-final`.
We will reuse the Kafka `docker-compose.yml` from Day 15.

#### **2. The Command Service (Write-Only)**

This service only accepts `POST` requests and writes to Kafka.
_Create `command/main.go`:_

```go
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
```

#### **3. The Query Service (Read-Only)**

This service runs a Kafka Consumer in the background to build its database, and exposes an HTTP `GET` endpoint for the frontend.
_Create `query/main.go`:_

```go
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
		Brokers:  []string{"localhost:9092"},
		GroupID:  "query-service-group",
		Topic:    "price-events",
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
```

---

### **Actionable Task for Today**

1. Make sure your Kafka container is running.
2. Open Terminal 1: Run the Query Service `go run query/main.go`.
3. Open Terminal 2: Run the Command Service `go run command/main.go`.
4. Open Terminal 3 (or your browser) and test the CQRS flow:
   - First, try to read: `curl http://localhost:8082/get-price?item=Nakroth` (Returns "Item not found").
   - Now, write data via the Command API: `curl -X POST "http://localhost:8081/update-price?item=Nakroth&price=15"`
   - Watch Terminal 1. You will see the Kafka consumer instantly pick it up and update the in-memory map.
   - Try reading again from the Query API: `curl http://localhost:8082/get-price?item=Nakroth` (Returns "$15").

---

### **End of Week 3 Review & Question**

You have survived the most conceptually difficult week of the roadmap. You now understand how large-scale data platforms like Uber, Netflix, and Amazon handle millions of events securely and quickly.

Take a massive breather! Tomorrow, we start **Week 4: Resilience & Distributed Transactions**. We will look at what happens when microservices go rogue and how to clean up the mess.

**To kick off Week 4, think about this scenario:**
You are booking a trip. You have a `Flight Service`, a `Hotel Service`, and a `Rental Car Service`.
You want to book all three, but it's an "all or nothing" deal. If you can't get the rental car, you want to cancel the flight and the hotel.

In an old monolithic SQL database, you would wrap all three inserts in a `BEGIN TRANSACTION` and `ROLLBACK` if one failed.
**Because these are now three entirely separate microservices with their own separate databases, how do you handle a "rollback" when the 3rd service fails?**
