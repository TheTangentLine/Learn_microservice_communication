### **Day 10: RabbitMQ Basics (Writing the Code)**

Today, we are going back to writing code. We will write a Producer (which sends a message) and a Consumer (which reads the message) using Go.

#### **1. The AMQP Protocol**

RabbitMQ uses a protocol called AMQP. To use it, you always follow these steps:

1. **Connect:** Open a TCP connection to the RabbitMQ server.
2. **Channel:** Open a lightweight "Channel" inside that TCP connection (this saves resources so you don't need a new TCP connection for every thread).
3. **Declare:** Declare the Queue. (You do this in _both_ the producer and the consumer because you never know which service will start up first!).
4. **Publish/Consume:** Send or read the data.

#### **2. Project Setup**

Inside your `week2-async` folder, initialize a Go module and install the official RabbitMQ package:

```bash
go mod init week2-async
go get github.com/rabbitmq/amqp091-go
```

Create two folders: `producer/` and `consumer/`.

#### **3. The Producer Code**

In `producer/main.go`, we will connect to the broker and send a single "OrderPlaced" message.

```go
package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. Connect to RabbitMQ (Running in our Docker container from Day 9)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 2. Open a Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 3. Declare a Queue to ensure it exists
	q, err := ch.QueueDeclare(
		"order_queue", // name
		false,         // durable (does it survive a broker restart?)
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// 4. Publish a message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "OrderPlaced: User123 bought Nakroth Skin"

	err = ch.PublishWithContext(ctx,
		"",     // exchange (we will learn this tomorrow)
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
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

In `consumer/main.go`, we will connect to the same queue and wait for messages.

```go
package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 2. Open a Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 3. Declare the exact same Queue (in case Consumer starts before Producer)
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

	// 4. Register a consumer
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer name
		true,   // auto-ack (important for tomorrow's question!)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// 5. Read messages forever using a Go channel
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
2. Open two terminals.
3. In Terminal 1, start the consumer: `go run consumer/main.go`. It will hang there, listening.
4. In Terminal 2, run the producer: `go run producer/main.go`.
5. Watch Terminal 1 instantly print out the message!
6. **Experiment:** Stop the consumer (CTRL+C). Run the producer 5 times. Open the RabbitMQ web UI (`http://localhost:15672`), click on "Queues", and look at the "order_queue". You will see 5 messages sitting there, waiting patiently. Start the consumer again, and watch it instantly drain all 5 messages. **This is temporal decoupling in action.**

---

### **Day 10 Revision Question**

Look closely at step 4 in the Consumer code. There is a parameter called `auto-ack` (auto-acknowledge), which we set to `true`.

This tells RabbitMQ: "As soon as I pull this message from the queue, delete it from the broker immediately."

**If our Consumer is a payment processing service, why is setting `auto-ack: true` an incredibly dangerous idea? What happens if our code panics/crashes halfway through processing the credit card?**

---

**Answer:**

### **The Danger of `auto-ack: true` (At-Most-Once Delivery)**

If your consumer is a Payment Service and `auto-ack` is true:

1. RabbitMQ hands the "Process Payment $50" message to the worker.
2. RabbitMQ _instantly_ deletes the message from the queue.
3. The worker starts processing the credit card... and then out of memory (OOM) crashes.
4. **The result:** The message is gone forever. The user is never charged, and the order is permanently stuck. This is a disaster.

### **The Fix: Manual ACKs (At-Least-Once Delivery)**

To fix this, we set `auto-ack: false`. The worker must explicitly tell RabbitMQ, _"I am completely done processing this, you can delete it now."_

1. RabbitMQ hands the message to the worker, but keeps a copy safely in the queue (marked as 'unacknowledged').
2. The worker successfully charges the credit card.
3. **CRASH!** Right before the worker can send the `ch.Ack()` back over the network, the server loses power.

### **The Idempotency Problem**

Here is where your intuition kicked in. Because RabbitMQ never received the ACK, it assumes the worker died _before_ processing the message.
RabbitMQ instantly re-queues the message and hands it to **Worker #2**.

If your worker is NOT idempotent, Worker #2 will charge the user's credit card _again_. You just double-charged your customer because of a network blip.

What you referred to as a "race condition" is more accurately called the **Duplicate Message Problem**. In distributed systems, because networks are unreliable (Fallacy #1!), you can never guarantee a message is delivered exactly once. You can only guarantee it is delivered _at least once_.

Therefore, **every consumer must be idempotent.** Idempotent means no matter how many times you apply an operation, the result is the same as applying it once (e.g., `x = 5` is idempotent; `x = x + 5` is not).

### **How to Make a Worker Idempotent**

You handle this in your database using **Idempotency Keys**:

1. Every message must have a unique ID (e.g., `order_id: 999`).
2. When the Payment Worker receives the message, the _very first thing_ it does is check its own database: `SELECT status FROM payments WHERE order_id = 999`.
3. If the status is `completed`, the worker says, _"Ah, I already did this!"_ It safely skips the credit card charge, sends the `ACK` to RabbitMQ to clear the duplicate message, and moves on.
