package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var q amqp.Queue
var ch *amqp.Channel

func generateUUID() string {
	id := uuid.New().String()
	return "order_" + id
}

func publishToRabbitMQ(event string) error {
	// 4. Publish a message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ch.PublishWithContext(ctx,
		"inventory", // exchange
		"",          // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		})
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return err
	}

	log.Printf(" [x] Sent %s\n", event)
	return nil
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	orderID := generateUUID() // e.g., "ord_12345"

	// 1. Create the JSON event
	event := fmt.Sprintf(`{"order_id": "%s", "item": "%s"}`, orderID, item)

	// 2. Publish to RabbitMQ (Fire and Forget)
	err := publishToRabbitMQ(event)
	if err != nil {
		http.Error(w, "Failed to place order", http.StatusInternalServerError)
		return
	}

	// 3. Immediately respond to the user
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf("Success! Order %s is processing.\n", orderID)))
}

func main() {
	// 1. Connect to RabbitMQ (Running in our Docker container from Day 9)
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 2. Open a Channel
	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 3. Declare a Queue to ensure it exists
	err = ch.ExchangeDeclare(
		"inventory", // name
		"fanout",    // type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	http.HandleFunc("/", checkoutHandler)
	if err = http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
