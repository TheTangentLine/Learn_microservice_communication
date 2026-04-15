### **Day 18: Redis Pub/Sub & In-Memory Messaging**

Today, we look at **Redis**. You probably know Redis as a blazing-fast, in-memory caching database. But it also has a built-in messaging system called **Redis Pub/Sub**.

If Kafka is a heavy, durable hard drive that remembers everything forever, Redis Pub/Sub is a fleeting whisper in the wind.

#### **1. How Redis Pub/Sub Works**

- **The Concept:** A Publisher sends a message to a "Channel" (like a Topic). Subscribers listen to that Channel.
- **The Catch (100% Ephemeral):** Redis does **not** store Pub/Sub messages. Not on disk, not even in memory. It literally just receives the message and immediately pushes it down the TCP socket to whoever is currently listening.
- **The Consequence:** If your Consumer Service crashes, and a message is published while it's rebooting... that message is gone forever. There is no queue to hold it. There is no offset to rewind.

#### **2. Why use it if it loses data?**

Because it is **unbelievably fast** and has incredibly low latency.
You use Redis Pub/Sub for data where _history doesn't matter_, and you only care about _right now_.

- **Good Use Cases:**
  - Live multiplayer game coordinates (if a player drops a packet, you don't care where they were 2 seconds ago, you just need their next coordinate).
  - Live streaming chat rooms (Twitch/YouTube chat).
  - Real-time stock price tickers on a UI.
- **Terrible Use Cases:**
  - Order processing.
  - Payments.
  - Inventory management.

#### **3. The Microservice "WebSocket" Pattern**

This is the most common use of Redis Pub/Sub in microservices.
Imagine you have 5 instances of a "WebSocket Service" keeping connections open to 100,000 user browsers.
Your `Order Service` finishes processing an order. It publishes a message to a Redis channel: `user_123_updates`.
All 5 WebSocket services hear it instantly, but only the specific worker holding the TCP connection for `user_123` forwards the message to the browser, updating their UI without a page refresh!

---

### **Actionable Task for Today**

Let's see this "whisper in the wind" in action using Docker.

1.  **Spin up Redis:** Open a terminal and run a quick Redis container:
    ```bash
    docker run --name my-redis -p 6379:6379 -d redis
    ```
2.  **Open Terminal 1 (The Subscriber):** Access the Redis CLI inside the container and subscribe to a channel:
    ```bash
    docker exec -it my-redis redis-cli
    > SUBSCRIBE live_chat
    ```
    _(It will sit there, waiting)._
3.  **Open Terminal 2 (The Publisher):**
    Open another CLI and publish a message:
    ```bash
    docker exec -it my-redis redis-cli
    > PUBLISH live_chat "Hello from the Publisher!"
    ```
4.  **Watch Terminal 1:** It instantly receives the message.
5.  **The Ephemeral Test:** In Terminal 1, hit `CTRL+C` to kill the subscriber. In Terminal 2, publish 3 more messages. Start Terminal 1 and `SUBSCRIBE live_chat` again.
    **Notice how those 3 messages are completely gone. Redis didn't save them for you.**

---

### **Day 18 Revision Question**

You are building a chat application like WhatsApp or Discord.

You need the absolute blazing real-time speed of **Redis Pub/Sub** when both User A and User B have the app open on their screens.
However, if User B's phone is turned off, you cannot afford to lose the message—it needs to be delivered when they turn their phone back on tomorrow (which Redis cannot do).

**How would you architect a system that combines the tools we've learned (Databases, Queues/Kafka, and Redis Pub/Sub) to get both real-time speed AND durable delivery?**

**Corrected Answer:**

### **The Redis Nuance: Storage vs. Pub/Sub**

You are absolutely right that Redis can save data to disk using RDB (Snapshots) or AOF (Append-Only Files). If you run `SET my_key "hello"`, Redis will save it to the hard drive.

**However, Redis Pub/Sub completely ignores this.** Even if you have disk persistence turned on, a `PUBLISH` command is a fire-and-forget operation in the Redis engine. It doesn't write the message to the AOF log, and it doesn't store it in memory. If a subscriber isn't connected at that exact millisecond, the message vanishes into the ether.

_(Note: In Redis 5.0, they released a feature called **Redis Streams** specifically to act like Kafka and solve this, but classic Pub/Sub remains ephemeral)._

---

### **Your Architecture: The "Dual-Write" Problem**

Your idea to publish to both—"1 to Redis, 1 to Kafka"—is the most logical first thought every engineer has when trying to get the best of both worlds.

But imagine the code in your Chat Service:

```go
err1 := publishToRedis(msg) // For instant real-time delivery
err2 := publishToKafka(msg) // For durable backup
```

**The Fatal Flaw:** What happens if `publishToRedis` succeeds, but the network blips for a microsecond and `publishToKafka` fails?
User A sees the message on their screen instantly (via Redis), but the message was never saved to Kafka. When User A refreshes their app, the message is completely gone.

This is known as the **Dual-Write Problem**. You can never guarantee that two separate database/broker writes will both succeed at the exact same time without complex, slow distributed transactions.

---

### **The Industry Solution: The Relay Pattern**

To solve this, we rely on the single source of truth (Kafka) and use a "Relay" to get the speed of Redis. Here is how Discord or WhatsApp handles it:

1. **The Write (Durable):** User A sends a message. The Chat Service writes it to **Kafka only**. (We know it's safe now).
2. **The Database Consumer:** Worker 1 reads from Kafka and saves the message permanently to PostgreSQL/Cassandra. (This is for chat history).
3. **The Redis Relay Consumer:** Worker 2 reads from the exact same Kafka topic and instantly `PUBLISH`es the message to **Redis Pub/Sub**.
4. **The WebSockets:** The WebSocket servers, listening to Redis, instantly push the message to User B's open phone.

**What if User B's phone is off?**
They miss the Redis Pub/Sub broadcast entirely. But that's okay! When they turn their phone on tomorrow, their app does a standard REST API `GET /messages` to the PostgreSQL database to fetch the history.

You have successfully separated the _Write_ logic from the _Read_ logic. Which leads us perfectly into today's lesson.

---
