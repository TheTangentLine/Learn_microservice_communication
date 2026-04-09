package main

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var idempotencyDb map[string]bool

func databaseContains(orderID string) bool {
	_, ok := idempotencyDb[orderID]
	return ok
}

func saveToDatabase(orderID string) {
	idempotencyDb[orderID] = true
}

func processMessage(msg amqp.Delivery) {
	// 1. Parse JSON
	var event map[string]string
	json.Unmarshal(msg.Body, &event)
	orderID := event["order_id"]

	// 2. Check Idempotency (Simulated Database Check)
	if databaseContains(orderID) {
		log.Printf("Duplicate order %s detected. Skipping.", orderID)
		msg.Ack(false) // Acknowledge to clear it from the queue
		return
	}

	// 3. Do the Work
	log.Printf("Processing inventory deduction for order: %s", orderID)
	saveToDatabase(orderID)

	// 4. Manual Acknowledgment
	msg.Ack(false)
	log.Printf("Successfully processed order: %s", orderID)
}

func main() {
	idempotencyDb = make(map[string]bool)

	// 1. Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
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

	// 3.  Declare the Exchange (same as producer)
	err = ch.ExchangeDeclare("inventory", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register an Exchange: %v", err)
	}
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a queue: %v", err)
	}
	err = ch.QueueBind(
		q.Name,      // queue name
		"",          // routing key
		"inventory", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind a queue: %v", err)
	}

	// 4. Register a consumer
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer name
		true,   // auto-ack
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
		for msg := range msgs {
			processMessage(msg)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
