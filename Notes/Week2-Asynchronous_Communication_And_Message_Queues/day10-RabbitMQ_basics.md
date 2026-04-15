### **Day 10: RabbitMQ Basics (Writing the Code)**

Today we write a Producer (sends a message) and a Consumer (reads a message) using Go and the AMQP protocol.

#### **1. The AMQP Connection Steps**

Every interaction with RabbitMQ follows this sequence:

```mermaid
sequenceDiagram
    participant Producer
    participant RabbitMQ
    participant Consumer

    Producer->>RabbitMQ: Connect (TCP)
    Producer->>RabbitMQ: Open Channel
    Producer->>RabbitMQ: QueueDeclare("order_queue")
    Producer->>RabbitMQ: Publish(message)
    RabbitMQ-->>Producer: OK

    Consumer->>RabbitMQ: Connect (TCP)
    Consumer->>RabbitMQ: Open Channel
    Consumer->>RabbitMQ: QueueDeclare("order_queue")
    Consumer->>RabbitMQ: Consume(auto-ack: false)
    RabbitMQ-->>Consumer: Deliver message (unacknowledged)
    Consumer->>Consumer: process message
    Consumer->>RabbitMQ: ACK (manual)
    RabbitMQ->>RabbitMQ: delete message from queue
```

Both the producer and consumer declare the same queue — because you never know which service starts up first.

#### **2. Project Setup**

```bash
go mod init week2-async
go get github.com/rabbitmq/amqp091-go
```

Create folders: `producer/` and `consumer/`.

#### **3. The Producer Code**

In `producer/main.go`:

```go
package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"order_queue", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "OrderPlaced: User123 bought Nakroth Skin"

	err = ch.PublishWithContext(ctx,
		"",     // exchange (covered tomorrow)
		q.Name, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf(" [x] Sent %s\n", body)
}
```

#### **4. The Consumer Code**

In `consumer/main.go`:

```go
package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"order_queue",
		false, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // auto-ack — see revision question below!
		false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
```

---

### **Actionable Task for Today**

1. Ensure your RabbitMQ Docker container from Day 9 is running (`docker ps`).
2. Terminal 1: `go run consumer/main.go` — it will hang, listening.
3. Terminal 2: `go run producer/main.go`.
4. Watch Terminal 1 instantly print the message.

**The key experiment:** Stop the consumer (CTRL+C). Run the producer 5 times. Open the RabbitMQ UI, click "Queues," and observe 5 messages sitting safely in `order_queue`. Start the consumer again — watch it drain all 5 instantly. **This is temporal decoupling in action.**

---

### **Day 10 Revision Question**

In the Consumer, `auto-ack` is set to `true`. This tells RabbitMQ: "Delete the message from the queue the moment I pull it."

**If our Consumer is a Payment Service, why is `auto-ack: true` dangerous? What happens if our code crashes halfway through charging the credit card?**

**Answer: The Danger of auto-ack (At-Most-Once Delivery)**

1. RabbitMQ delivers the "Process Payment $50" message to the worker.
2. RabbitMQ **instantly deletes** the message from the queue.
3. The worker starts charging the credit card... and crashes (OOM, server dies).
4. **Result:** The message is gone forever. The user is never charged, the order is permanently stuck. A disaster.

**The Fix: Manual ACKs (At-Least-Once Delivery)**

Set `auto-ack: false`. The worker must explicitly tell RabbitMQ "I am completely done — you can delete it now."

```mermaid
flowchart TD
    RMQ["RabbitMQ"]
    Worker["Payment Worker"]

    subgraph autoack ["auto-ack: true (DANGEROUS)"]
        direction LR
        A1["RabbitMQ delivers message"] --> A2["RabbitMQ deletes immediately"]
        A2 --> A3["Worker crashes mid-charge"]
        A3 --> A4["Message lost forever"]
    end

    subgraph manualack ["auto-ack: false (SAFE)"]
        direction LR
        M1["RabbitMQ delivers + keeps copy"] --> M2["Worker charges card successfully"]
        M2 --> M3["Worker sends ch.Ack()"]
        M3 --> M4["RabbitMQ deletes message"]
        M1 --> M5["Worker crashes before Ack"]
        M5 --> M6["RabbitMQ re-queues for another worker"]
    end
```

**The Idempotency Problem:** Because RabbitMQ re-queues on crash, a second worker might process the same message and charge the credit card again. This is the **Duplicate Message Problem**. The solution: every consumer must be **idempotent** (covered in detail on Day 11 and Day 13).
