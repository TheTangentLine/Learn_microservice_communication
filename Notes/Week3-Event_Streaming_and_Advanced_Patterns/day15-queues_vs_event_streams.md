### **Week 3: Event Streaming & Advanced Patterns**

### **Day 15: Queues vs. Event Streams (The Paradigm Shift)**

Today, we are shifting from **Message Queues** (RabbitMQ/SQS) to **Event Streams** (Apache Kafka). To understand Kafka, you have to unlearn how RabbitMQ works.

#### **1. The Queue Mindset (Transient)**

- **Analogy:** A To-Do List.
- **How it works:** You write down a task. A worker reads it, does the job, and crosses it off (deletes it).
- **The Problem:** Once it's crossed off, it is gone forever. If a new worker joins tomorrow and asks, "What tasks did we do yesterday?", the queue has no idea. The queue is completely empty.

#### **2. The Stream Mindset (Durable & Append-Only)**

- **Analogy:** A Captain's Log or a Diary.
- **How it works:** Events are written sequentially into an **Append-Only Log** on a hard drive. You can add to the bottom, but you can _never_ delete or modify what is already written.
- **The Magic:** Because the broker doesn't delete the messages, the broker doesn't care who reads them. Instead, the **Consumers** are responsible for keeping track of their own place in the book using a bookmark called an **Offset**.

#### **3. Smart Broker vs. Dumb Broker**

- **RabbitMQ is a "Smart Broker / Dumb Consumer."** RabbitMQ does all the heavy lifting. It routes messages, tracks exactly which worker has what, waits for acknowledgments, and deletes the data.
- **Kafka is a "Dumb Broker / Smart Consumer."** Kafka literally just dumps bytes onto a hard drive as fast as possible. The Consumers have to be smart enough to remember their own offsets and pull the data themselves. Because Kafka does so little work, it can process _millions_ of messages per second, vastly outperforming RabbitMQ in sheer throughput.

---

### **Actionable Task for Today**

You will need Docker again, but Kafka is a heavier beast than RabbitMQ.

**1. Create a new folder:** `week3-streaming`
**2. Create a `docker-compose.yml` file:**
We will use a modern Kraft-based Kafka image (which removes the need for Zookeeper, an older dependency you might see in older tutorials).

```yaml
version: "3.8"
services:
  kafka:
    image: bitnami/kafka:latest
    ports:
      - "9092:9092"
    environment:
      - KAFKA_ENABLE_KRAFT=yes
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092
      - KAFKA_BROKER_ID=1
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@kafka:9093
      - ALLOW_PLAINTEXT_LISTENER=yes
```

Run `docker-compose up -d` to get it downloading and running. Tomorrow, we will connect to it!

---

### **Day 15 Revision Question**

Think about the "Event Stream Mode" in the widget above, where messages are permanently stored on the hard drive.

If Amazon processes 50 million orders a day, and Kafka never deletes a message when it is read, **what is going to eventually happen to the Kafka server, and how do you think Kafka solves this practical reality?**

**Answer:**

"Sliding window" is the exact architectural concept!

If Kafka literally never deleted anything, the hard drives would eventually fill up and the server would crash. To prevent this, Kafka implements a sliding window **Retention Policies**.

Kafka handles this in three ways:

1.  **Time-Based Retention (The Default):** Kafka's default sliding window is **7 days**. As soon as a message is 7 days and 1 second old, Kafka quietly drops it off the back of the log to free up space.
2.  **Size-Based Retention:** You can tell Kafka, _"I don't care about time, just make sure this topic never exceeds 500GB."_
3.  **Log Compaction:** This is a special Kafka feature. If your topic is a stream of database updates, Kafka can look at the keys (e.g., `user_123`) and say, _"I have 50 updates for user_123 in this log. I'm going to delete the 49 old ones and only keep the newest one."_ You have the intuition of a senior data engineer. Let's get into the mechanics.
