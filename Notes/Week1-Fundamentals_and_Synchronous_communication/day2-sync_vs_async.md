### **Day 2: Sync vs. Async Overview**

Today is about understanding the two fundamental ways microservices communicate. You don't need to write code today; you need to understand the architectural trade-offs so you know _which_ code to write tomorrow.

#### **1. Synchronous Communication (Request/Response)**

Think of this like a **phone call**. You ask a question, and you wait on the line until you get an answer.

- **How it works:** Service A sends a request (usually HTTP/REST or gRPC) to Service B and blocks (waits) until Service B sends a response back.
- **The Good:** It's conceptually simple. You know immediately if the action succeeded or failed.
- **The Bad:** It creates **temporal coupling**. Both services must be alive and healthy at the exact same time. If Service B is slow, Service A becomes slow (cascading failure).
- **When to use it:** When you _absolutely must_ have an immediate answer before proceeding (e.g., verifying a user's password during login).

#### **2. Asynchronous Communication (Event-Driven / Message Passing)**

Think of this like an **email** or a **text message**. You send it, go about your day, and the receiver deals with it when they are ready.

- **How it works:** Service A drops a message into a middleman (a Message Broker or Queue like RabbitMQ or Kafka) and immediately moves on. Service B picks up the message from the broker at its own pace.
- **The Good:** Total decoupling. If Service B goes down for 10 minutes, Service A can keep accepting traffic and dropping messages into the queue. When Service B comes back online, it just processes the backlog. It handles massive spikes in traffic beautifully.
- **The Bad:** It introduces **eventual consistency**. Service A doesn't know if Service B successfully processed the task. Error handling gets much more complex.
- **When to use it:** Background tasks, notifications, or processes that take a long time (e.g., generating a PDF report, sending an email, or resizing an image).

---

### **Actionable Task: Mental Mapping**

Grab a piece of paper or a digital whiteboard. Map out a **User Registration Flow** using both paradigms.

1.  **Sync Approach:** The `Gateway` calls `User Service` -> calls `Email Service` to send a welcome email -> responds to the user. What happens to the user's request if the `Email Service` is down?
2.  **Async Approach:** The `Gateway` calls `User Service` -> saves user to DB and pushes a "UserCreated" event to a queue -> responds to the user immediately. The `Email Service` reads the queue and sends the email.

_Hint for tomorrow:_ When we start writing code on Day 3, we will be building the Sync approach first using Go or Python.

---

### **Day 2 Revision Question**

Imagine you are building an e-commerce checkout flow. When a user clicks "Place Order," two things must happen:

1.  Charge their credit card via a third-party banking API.
2.  Update the inventory system to deduct the items they bought.

Which of these two steps should ideally be **Synchronous**, and which should be **Asynchronous**? Why?

**Answer:**

1. **Payment:** You generally want the _authorization_ to be **Synchronous**. You need to know right then and there if the card is valid and has funds before the user leaves the checkout screen. (Though actually capturing the funds later can be async).
2. **Inventory:** A queue (**Asynchronous**) is the absolute best way to prevent race conditions. If 1,000 people try to buy the last 5 items at the exact same millisecond, a queue perfectly lines up those requests one by one so you don't accidentally oversell. Amazon relies heavily on this "eventual consistency."
