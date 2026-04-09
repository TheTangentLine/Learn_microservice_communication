### **Day 11: Advanced Routing (Exchanges)**

Yesterday, we published directly to a Queue. But in real Publish/Subscribe (Pub/Sub) architectures, Producers _never_ send messages directly to queues.

#### **1. The Exchange**

Producers only send messages to an **Exchange**. The Exchange is like a mail sorting facility. It looks at the message and decides which Queue(s) it should be copied into, based on rules called **Bindings**.

#### **2. Types of Exchanges**

- **Fanout Exchange:** The simplest. It blindly copies the message to _every_ queue bound to it. (Perfect for our `OrderPlaced` event going to both Inventory and Email queues).
- **Direct Exchange:** Routes messages based on an exact matching word (Routing Key). E.g., Send messages with key `pdf_tasks` to the PDF Queue, and `image_tasks` to the Image Queue.
- **Topic Exchange:** The most powerful. It routes based on wildcard patterns.
  - E.g., A routing key might be `order.europe.shoes`.
  - Queue A listens for `order.*.shoes` (All shoe orders globally).
  - Queue B listens for `order.europe.#` (All European orders of any kind).

---

### **Actionable Task for Today**

Let's modify yesterday's code to use a **Fanout Exchange** so we can have multiple _different_ services listen to the same event.

**1. Update the Producer (`producer/main.go`):**
Remove the `QueueDeclare` block entirely. Instead, declare an Exchange:

```go
// Declare a Fanout Exchange named "logs_exchange"
err = ch.ExchangeDeclare(
	"logs_exchange", // name
	"fanout",        // type
	true,            // durable
	false,           // auto-deleted
	false,           // internal
	false,           // no-wait
	nil,             // arguments
)
```

When publishing, change the parameters to publish to the exchange, not a specific queue:

```go
err = ch.PublishWithContext(ctx,
	"logs_exchange", // Publish to our new Exchange!
	"",              // Routing key is ignored for fanout
	// ... rest remains the same
)
```

**2. Update the Consumer (`consumer/main.go`):**
The consumer needs to create a temporary, uniquely named queue and _bind_ it to the Exchange.

```go
// 1. Declare the Exchange (same as producer)
err = ch.ExchangeDeclare("logs_exchange", "fanout", true, false, false, false, nil)

// 2. Declare a temporary, exclusive queue (RabbitMQ generates a random name)
q, err := ch.QueueDeclare("", false, false, true, false, nil)

// 3. BIND the queue to the exchange
err = ch.QueueBind(
	q.Name,          // queue name
	"",              // routing key
	"logs_exchange", // exchange
	false,
	nil,
)
```

**Run the test:**

1. Open _three_ terminals.
2. Run the Consumer in Terminal 1 and Terminal 2. (These simulate your Inventory Service and Email Service).
3. Run the Producer in Terminal 3.
4. Watch the exact same message perfectly fan out and appear in both Consumer terminals simultaneously!

---

### **Day 11 Revision Question**

You’ve got a solid grasp on how to fix the "duplicate message" problem using an Idempotency Key in your database.

But think about the database itself. If 5 duplicate messages hit your worker at the exact same millisecond, and all 5 concurrent threads run `SELECT status FROM payments WHERE order_id = 999` before any of them have a chance to write `status = 'completed'`... all 5 might think the order hasn't been processed yet!

**Answer:**

**How do you solve this specific database-level race condition to ensure your idempotency check is actually bulletproof?** (Hint: Think about primary keys or specific SQL constraints).

Here is the industry-standard solution:

1. Create a specific table in your database called `idempotency_keys` with a single column: `key_id` (String).
2. Make `key_id` the **Primary Key** (or give it a `UNIQUE` constraint).
3. When your Payment Worker gets a message, it doesn't run a `SELECT`. Instead, it immediately tries to run an `INSERT INTO idempotency_keys (key_id) VALUES ('payment_order_999')`.

**What happens to the 5 concurrent threads?**

- **Thread 1** hits the database first. The insert succeeds. Thread 1 moves on to charge the credit card.
- **Threads 2, 3, 4, and 5** hit the database a millisecond later. The database strictly enforces the Primary Key rule and instantly throws a **"Unique Constraint Violation" (Duplicate Key Error)** for all four threads.
- Your Go code catches that specific error, says _"Ah, another thread is already handling this!"_, safely ignores the credit card charge, and sends the `ACK` to RabbitMQ.

Bulletproof idempotency achieved!
