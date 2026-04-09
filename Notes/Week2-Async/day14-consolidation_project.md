### **Day 14: Week 2 Consolidation Project**

Today, we are tearing down our Synchronous Week 1 project and rebuilding it as a highly resilient, Event-Driven Architecture.

#### **The Architecture**

1. **User** sends an HTTP POST request to the API Gateway.
2. **API Gateway** routes it to the **Order Service**.
3. **Order Service** instantly returns an HTTP 202 (Accepted) to the user: _"Your order is being processed!"_ 4. **Order Service** (Producer) publishes an `OrderPlaced` event to a RabbitMQ Fanout Exchange.
4. **Inventory Service** (Consumer) listens to the queue, enforces idempotency, deducts the stock, and sends a manual `ACK` to RabbitMQ.

#### **1. Project Setup**

Create a new folder: `week2-final`.
Your structure should look like this:

```text
week2-final/
├── docker-compose.yml
├── gateway/          # API Gateway (HTTP proxy)
├── order/            # HTTP Server + RabbitMQ Producer
└── inventory/        # RabbitMQ Consumer with Idempotency logic
```

#### **2. The Docker Compose File**

We need to spin up our services alongside RabbitMQ.

```yaml
version: "3.8"
services:
  rabbitmq:
    image: rabbitmq:3-management-alpine
    ports:
      - "5672:5672"
      - "15672:15672"

  inventory:
    build: ./inventory
    depends_on:
      - rabbitmq

  order:
    build: ./order
    depends_on:
      - rabbitmq

  gateway:
    build: ./gateway
    ports:
      - "8000:8000"
    depends_on:
      - order
```

#### **3. The Order Service (The Producer)**

This service no longer waits for the Inventory service. It just drops a message and responds instantly.

```go
// order/main.go (Simplified Snippet)
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
```

#### **4. The Inventory Service (The Consumer)**

This service runs in the background, carefully processing messages and acknowledging them.

```go
// inventory/main.go (Simplified Snippet)
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
```

---

### **Actionable Task for Today**

1. Build out the folder structure and write the Docker files (you can reuse the Dockerfiles and Gateway code from Week 1!).
2. Write the Go code using the AMQP library we learned on Day 10.
3. Run `docker-compose up --build`.
4. Hit your gateway: `http://localhost:8000/api/checkout?item=Nakroth`
5. Watch the logs. You will see the Order service respond to you instantly, while a fraction of a second later, the Inventory service logs that it picked up the message and processed it.

---

### **End of Week 2 Review & Question**

You have officially conquered Message Queues, AMQP, Fanouts, Visibility Timeouts, and Idempotency. This is the bread and butter of senior backend engineering. Take a massive breather, you earned it.

Tomorrow, we start **Week 3: Event Streaming & Advanced Patterns**, where we introduce the undisputed heavyweight champion of distributed data: **Apache Kafka**.

**To wrap up Week 2:**
RabbitMQ deletes messages the exact moment they are successfully acknowledged by a consumer. The queue stays empty.
If we want a new Data Analytics service to calculate "Total Items Sold Today", but the messages have already been deleted by the Inventory service hours ago, RabbitMQ can't help us.

**Without looking ahead to Kafka, how might an architectural system solve the problem of needing to read historical events that have already happened and been processed?**
