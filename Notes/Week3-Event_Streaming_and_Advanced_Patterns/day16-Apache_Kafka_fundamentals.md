### **Day 16: Apache Kafka Fundamentals**

Today we map your "Queue" vocabulary to your new "Stream" vocabulary.

#### **1. The Core Vocabulary**

- **Topic (Instead of Queue):** A named stream of records. Think of it like a database table. You don't publish to a queue; you publish to a Topic (e.g., `order-events`).
- **Partition (The Secret to Scale):** A single topic is split into multiple Partitions (Log 0, Log 1, Log 2). They are spread across different physical hard drives. This means if 10,000 producers are writing to the same topic simultaneously, they aren't fighting over one file lock; they are writing to multiple partitions in parallel.
- **Offset:** A sequential ID number assigned to a message as it arrives in a partition. (Message 0, Message 1, Message 2). This is the "bookmark" the consumer uses to know where it left off.
- **Consumer Group:** This is how Kafka scales reading. If you put 3 instances of your `Inventory Service` into the same "Consumer Group," Kafka will automatically divide the topic's partitions among them so they don't process duplicate data.

#### **2. How Kafka Routes Data**

In RabbitMQ, we used Exchanges to route data. Kafka doesn't have Exchanges.
When a Producer sends a message to Kafka, it can include a **Key** (like the `user_id`). Kafka hashes that Key and uses it to choose the Partition.

- _Rule of thumb:_ Messages with the _same key_ (e.g., all events for `user_123`) always go to the _same partition_. This guarantees they are processed in the exact order they occurred!

---

### **Actionable Task for Today**

We are going to use the Docker container you spun up yesterday and interact with it using Kafka's built-in CLI tools.

1. Ensure your Kafka container is running (`docker ps`).
2. **Open Terminal 1 (Create a Topic):**
   Run this command to execute a script _inside_ your Docker container. We are creating a topic with 3 partitions.

   ```bash
   docker exec -it <your-kafka-container-name> /opt/bitnami/kafka/bin/kafka-topics.sh \
     --create --topic order-events \
     --partitions 3 \
     --replication-factor 1 \
     --bootstrap-server localhost:9092
   ```

3. **Open Terminal 2 (The Producer):**
   Start the interactive producer CLI.

   ```bash
   docker exec -it <your-kafka-container-name> /opt/bitnami/kafka/bin/kafka-console-producer.sh \
     --topic order-events \
     --bootstrap-server localhost:9092
   ```

   _Your terminal will wait for input. Type `Hello Kafka` and press Enter. Type `Order 2 placed` and press Enter._

4. **Open Terminal 3 (The Consumer):**
   Start the interactive consumer CLI. We will tell it to read from the very beginning of the sliding window!
   ```bash
   docker exec -it <your-kafka-container-name> /opt/bitnami/kafka/bin/kafka-console-consumer.sh \
     --topic order-events \
     --from-beginning \
     --bootstrap-server localhost:9092
   ```

Watch how the consumer instantly prints out the messages you just typed.
_Experiment:_ Stop the consumer (CTRL+C). Go to Terminal 2 and type 5 more messages. Then restart the consumer in Terminal 3, **WITHOUT** the `--from-beginning` flag. What happens? Then stop it, and restart it **WITH** the `--from-beginning` flag.

You will clearly see that the messages were never deleted!

---

### **Day 16 Revision Question**

Kafka's **Consumer Groups** are the key to its massive scalability, but they have one strict, mathematical rule regarding **Partitions**.

Imagine you have a topic called `website-clicks` that has exactly **3 Partitions**.
You start up a Consumer Group for your Analytics Service, and because traffic is heavy, you spin up **5 instances (workers)** of your Analytics Service in that group.

**Based on how Kafka divides work, what exactly is going to happen to those 5 workers, and why is this an important architectural constraint to remember when designing your Kafka topics?**

**Answer:**

Here is the Golden Rule of Kafka: **One Partition can be consumed by AT MOST ONE Consumer within a specific Consumer Group.**

If you have 3 partitions and 5 workers in the same group:

- Worker 1 gets Partition 0.
- Worker 2 gets Partition 1.
- Worker 3 gets Partition 2.
- **Worker 4 and Worker 5 get NOTHING. They sit 100% idle.**

**Why does Kafka enforce this strict limitation?**
It is entirely about **Ordering Guarantees**. If Worker 4 and Worker 1 were both allowed to read from Partition 0 at the same time, they would race each other. Worker 4 might process the `DeleteUser` event before Worker 1 finishes processing the `CreateUser` event, destroying your database integrity.

**The Architectural Lesson:** Your maximum scale-out capacity is strictly limited by your partition count. If you think you might need 10 workers in the future, you _must_ create the topic with at least 10 partitions on Day 1!
