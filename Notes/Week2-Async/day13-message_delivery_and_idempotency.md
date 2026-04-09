### **Day 13: Message Delivery Guarantees & Idempotency Deep Dive**

Today, we formalize the chaotic nature of networks. In distributed systems, engineers categorize message delivery into three strict guarantees.

#### **1. The Three Delivery Guarantees**

- **At-Most-Once (The "Fire and Forget" approach):** \* _How it works:_ The message is sent. If it gets lost, oh well. (This is RabbitMQ with `auto-ack: true`).
  - _When to use it:_ IoT sensor data (if you drop one temperature reading, the next one will arrive in 5 seconds anyway), or logging non-critical analytics.
- **At-Least-Once (The Industry Standard):**
  - _How it works:_ The system guarantees the message will arrive, but due to network retries, it might arrive 2, 3, or 10 times. (This is RabbitMQ with Manual ACKs, and AWS SQS).
  - _When to use it:_ Almost everything in microservices. Financial transactions, emails, inventory updates. **This requires Idempotent consumers.**
- **Exactly-Once (The Holy Grail):**
  - _How it works:_ The message is delivered and processed one time, no matter what crashes.
  - _The truth:_ Mathematically, true "exactly-once" delivery over an unreliable network is considered impossible (look up the "Two Generals' Problem"). However, some systems like Apache Kafka simulate it using complex internal transactions, but it adds massive overhead.

#### **2. The Anatomy of an Idempotent Consumer**

We talked about the database `UNIQUE` constraint yesterday. Let's look at how that actually looks in Go pseudo-code. Every time you write an async worker, it should follow this exact pattern:

```go
func processPaymentMessage(msg []byte) {
    // 1. Extract the unique Idempotency Key from the message
    orderID := extractOrderID(msg)

    // 2. Try to insert this key into the DB to claim the lock
    err := db.Exec("INSERT INTO processed_events (event_id) VALUES (?)", orderID)

    if err != nil {
        if isUniqueConstraintViolation(err) {
            // WE ARE IDEMPOTENT! Another worker already did this.
            log.Println("Duplicate message detected. Skipping.")
            acknowledgeMessage(msg) // Tell the broker to delete it
            return
        }
        // It's a real database error (e.g., DB is down), do NOT acknowledge!
        // Let the message go back to the queue to be retried later.
        return
    }

    // 3. We successfully claimed the lock. Now do the actual work.
    chargeCreditCard()

    // 4. Acknowledge the message so the broker deletes it.
    acknowledgeMessage(msg)
}
```

---

### **Actionable Task for Today**

No new Docker containers today. Instead, review the code snippet above and map it to your Week 1 Consolidation Project.

Think about how you will inject this pattern into the `Inventory Service` for tomorrow. Tomorrow, we build the Week 2 Final Project, refactoring our Synchronous Week 1 system into a fully Asynchronous Event-Driven system using RabbitMQ!

---

### **Day 13 Revision Question**

Look closely at the `processPaymentMessage` Go code above.

Imagine this exact sequence of events happens:

1. The `INSERT` succeeds (Step 2).
2. The code moves to Step 3 and tries to `chargeCreditCard()`.
3. The 3rd party Stripe/Visa API is temporarily down, and `chargeCreditCard()` throws a fatal error and crashes our Go function before we reach Step 4.

Because we never reached Step 4, RabbitMQ will eventually re-queue the message and hand it to another worker. **When that second worker receives the message and starts at Step 1, what is going to happen, and what terrible state is our system now stuck in?**

**Answer:**

Instead of just checking if the ID _exists_, you use a **State Machine**.
Your `idempotency_keys` table needs a `status` column (`PENDING`, `COMPLETED`, `FAILED`) and a `created_at` timestamp.

1. Worker 1 inserts: `id=999, status=PENDING, time=12:00:00`.
2. Worker 1 crashes.
3. Worker 2 gets the message at 12:05:00. It queries the DB: `SELECT status, time FROM idempotency_keys WHERE id=999`.
4. Worker 2 sees `PENDING`, but notices the timestamp is 5 minutes old! It realizes Worker 1 must be dead.
5. Worker 2 updates the row, takes over the lock, charges the card, and updates the status to `COMPLETED`.
