### **Day 17: Kafka in Practice (Writing the Code)**

Today, we write a Go Producer and Consumer to connect to the Kafka broker we spun up yesterday. We will use the highly popular `kafka-go` library by Segment.

#### **1. Project Setup**

Inside your `week3-streaming` folder, initialize a Go module:

```bash
go mod init week3-streaming
go get github.com/segmentio/kafka-go
```

Create two folders: `producer/` and `consumer/`.

#### **2. The Producer Code**

Unlike RabbitMQ where we just connected and fired, the Kafka producer manages batches and connection pooling automatically in the background.

In `producer/main.go`:

```go
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
```

#### **3. The Consumer Code**

This is where the Consumer Group magic happens.

In `consumer/main.go`:

```go
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
		Brokers: []string{"localhost:9092"},
		GroupID: "analytics-group", // THIS IS CRITICAL!
		Topic:   "order-events",
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
```

---

### **Actionable Task for Today**

We are going to prove the Golden Rule of Kafka that we just discussed.

1. Make sure your Kafka Docker container is running.
2. Open **three** separate terminal windows.
3. In Terminals 1 and 2, run the Consumer: `go run consumer/main.go`.
   _(Because both have `GroupID: "analytics-group"`, Kafka will assign them different partitions from our 3-partition topic)._
4. In Terminal 3, run the Producer: `go run producer/main.go`.
5. **Watch closely:** Because we hardcoded the Key to `"user_99"`, the `Hash` balancer will hash that key and pick exactly _one_ partition. You will see that **only one** of your consumers does 100% of the work, and the other sits totally idle!
6. **Experiment:** Change the Go Producer code so the key is random (e.g., `"user_1"`, `"user_2"`). Run it again. Now you will see Kafka evenly distribute the work across both Consumer terminals!

---

### **Day 17 Revision Question**

We just saw how Kafka guarantees strict ordering by forcing all messages with the same Key into the same Partition, ensuring only one worker reads them.

But think about how HTTP gateways work. If 10,000 users hit our API Gateway at once, our `Order Service` spins up 10,000 concurrent Goroutines to handle those HTTP requests. Those Goroutines all try to call `w.WriteMessages` to Kafka simultaneously.

**Even if we use the exact same Key (`"user_99"`), is the _Order Service_ guaranteed to send those messages to Kafka in the exact order the user clicked the button? Why or why not?** Let me know what you think before we jump into Redis on Day 18!

**Answer:**

Even though Kafka strictly orders messages inside a single partition, **Kafka can only order messages based on when it _receives_ them, not when the user _clicked_ them.**

If 10,000 users click a button, their internet speeds are different. Even if they hit the API Gateway at the same time, Go's CPU scheduler will execute those 10,000 Goroutines in a completely unpredictable, random order. User #5 might hit the Kafka producer before User #1.

**The Golden Rule:** Kafka guarantees the order of _receipt_, not the order of _creation_. (If you absolutely need creation order, you have to attach client-side timestamps and sort them later in your analytics database!).
