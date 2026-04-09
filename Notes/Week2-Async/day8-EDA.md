### **Week 2: Asynchronous Communication & Message Queues**

### **Day 8: Intro to Event-Driven Architecture (EDA)**

Today, we are shifting our entire mental model. In Week 1, our services were bossy. The Order Service told the Inventory Service, _"Check this stock right now and give me the answer."_ In Week 2, our services become gossips. The Order Service will simply announce to the world, _"Hey everyone, an order was just placed!"_ and it doesn't care who is listening or what they do with that information.

#### **1. The Core Concepts of EDA**

Event-Driven Architecture is built on two massive benefits:

- **Temporal Decoupling:** As you noted, the services don't need to be alive at the same time. If the Inventory Service goes down for 5 minutes, the Order Service keeps taking orders and dropping messages into the queue. The queue holds them safely until the Inventory Service wakes up and processes the backlog.
- **Spatial Decoupling:** The Order Service no longer needs to know the IP address or DNS name of the Inventory Service (like we hardcoded in Week 1). It only needs to know where the Message Broker is.

#### **2. Commands vs. Events**

This is the most common stumbling block for developers learning async systems. You must understand the difference in language and intent:

- **A Command (What we did in Week 1):** \* **Naming:** Imperative verbs (`UpdateInventory`, `ChargeCreditCard`).
  - **Intent:** You want a specific action to happen, and you usually expect a response or confirmation. If it fails, the sender cares.
- **An Event (What we are doing now):**
  - **Naming:** Past-tense verbs (`OrderPlaced`, `PaymentSucceeded`, `UserRegistered`).
  - **Intent:** It is an immutable fact. Something _already happened_ in the past. The sender (Producer) publishes this fact to the broker and immediately forgets about it. The receivers (Consumers) decide if they care about that fact.

#### **3. The Trade-off: Eventual Consistency**

We trade the risk of the system crashing for the reality of **Eventual Consistency**.
When a user clicks "Buy," the Order Service instantly returns "Success! Your order is processing." But the inventory hasn't actually been deducted yet. It _will_ be deducted eventually (usually within milliseconds), but for a brief moment, your database is technically out of sync with reality.

---

### **Actionable Task for Today**

No coding today. Let's design the architecture we will build over the next few days.
Grab your notepad and map out a **Publish/Subscribe (Pub/Sub)** flow:

1.  **The Producer:** The User hits the API Gateway -> Order Service. The Order Service saves a "Pending" order to its own database.
2.  **The Event:** The Order Service publishes an `OrderPlaced` event (containing the `item_name` and `user_id`) to a Message Broker.
3.  **The Consumers:** Draw **two** separate services listening to the broker:
    - `Inventory Service`: Hears the event, deducts stock.
    - `Email Service`: Hears the _exact same event_, sends a "Thank you for your order" email to the user.

Notice how we added an Email Service without having to touch or update the Order Service code at all! That is the magic of EDA.

---

### **Day 8 Revision Question**

Eventual consistency creates interesting business challenges.

Imagine there is only **1** Nakroth skin left in stock.
User A and User B both click "Buy" at the exact same moment.
Because our Order Service is now fully asynchronous, it instantly accepts _both_ orders, publishes two `OrderPlaced` events to the queue, and tells both users "Success!"

When the Inventory Service reads those events from the queue a millisecond later, it realizes it only has 1 skin for 2 buyers. **In a real-world asynchronous system like Amazon, how do you handle this scenario where you accidentally accepted an order you can't fulfill?** Let me know your thoughts, and we will move on to Day 9 and spin up RabbitMQ in Docker!

**Answer:**
**The "Amazon Way" (Compensating Transactions)**
In the real world, Amazon takes a business-first approach. They don't make you wait on a loading spinner. They instantly tell you "Success!" and close the connection.
If the Inventory Service later realizes the item is gone, it publishes a new event: `OrderFailed`.
Other services listen to this to execute a **Compensating Transaction**:

1. The `Payment Service` catches it and issues a refund.
2. The `Email Service` catches it and sends you an email: _"We're so sorry, but the item you ordered is out of stock. We have refunded your card."_

It's better to apologize later than to make 100,000 users stare at a loading screen!
