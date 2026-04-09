### **Week 1: Fundamentals & Synchronous Communication**

_Goal: Understand the baseline of distributed systems and how services talk to each other in real-time._

- **Day 1: The Microservices Paradigm.** \* _Study:_ Monoliths vs. Microservices. Why split things up? Understand the challenges of distributed computing (network latency, partial failures).
  - _Action:_ Read about the "8 Fallacies of Distributed Computing."
- **Day 2: Sync vs. Async Overview.**
  - _Study:_ The difference between Request/Response (Synchronous) and Fire-and-Forget/Publish-Subscribe (Asynchronous).
  - _Action:_ Draw a diagram of a user checkout flow using both methods.
- **Day 3: RESTful Communication.**
  - _Study:_ HTTP/1.1, REST principles, JSON. The standard way services talk.
  - _Action:_ Build two simple local services (e.g., in Node.js, Python, or Go) where Service A calls Service B via HTTP.
- **Day 4: RPC & gRPC Fundamentals.**
  - _Study:_ What is Remote Procedure Call? Why use gRPC over REST? (HTTP/2, binary payloads, speed).
  - _Action:_ Read up on Protocol Buffers (Protobufs).
- **Day 5: Implementing gRPC.**
  - _Study:_ How gRPC streams data (unary, client streaming, server streaming, bidirectional).
  - _Action:_ Rewrite your Day 3 services to communicate using gRPC and a simple `.proto` file.
- **Day 6: API Gateways.**
  - _Study:_ The API Gateway pattern. How clients talk to microservices (routing, rate limiting, authentication).
  - _Action:_ Spin up an API Gateway locally (e.g., Kong, KrakenD, or AWS API Gateway docs) and route requests to your services.
- **Day 7: Week 1 Consolidation Project.**
  - _Action:_ Build a 3-tier sync architecture. An API Gateway routes to an `Order Service`, which synchronously calls an `Inventory Service` via REST or gRPC to check stock.

---

### **Week 2: Asynchronous Communication & Message Queues**

_Goal: Decouple your services so they don't have to wait for each other to finish tasks._

- **Day 8: Intro to Event-Driven Architecture (EDA).**
  - _Study:_ What are events? Commands vs. Events. Why decoupling increases system resilience.
- **Day 9: Message Queues & Brokers.**
  - _Study:_ What is a message broker? Point-to-point queues vs. Publish/Subscribe (Pub/Sub).
  - _Action:_ Install Docker if you haven't, as you'll need it for brokers.
- **Day 10: RabbitMQ Basics.**
  - _Study:_ Advanced Message Queuing Protocol (AMQP). Producers, Consumers, Queues.
  - _Action:_ Spin up a RabbitMQ Docker container. Write a script to push a message to a queue and another to consume it.
- **Day 11: RabbitMQ Advanced Routing.**
  - _Study:_ Exchanges (Direct, Topic, Fanout) and Bindings.
  - _Action:_ Create a "Fanout" exchange where one message from Service A is received by both Service B and Service C simultaneously.
- **Day 12: Cloud-Native Queues (AWS SQS/SNS).**
  - _Study:_ How managed cloud services differ from self-hosted. SQS (Queues) vs. SNS (Topics/Pub-Sub).
  - _Action:_ Look into the "SNS to SQS fanout pattern"—a very common real-world architecture.
- **Day 13: Idempotency & Message Delivery Guarantees.**
  - _Study:_ At-most-once, At-least-once, and Exactly-once delivery. Why consumers _must_ be idempotent (safe to process the same message twice).
- **Day 14: Week 2 Consolidation Project.**
  - _Action:_ Refactor your Week 1 project. When `Order Service` creates an order, it pushes an `OrderPlaced` event to RabbitMQ. `Inventory Service` listens to this queue, updates stock, and sends an `InventoryUpdated` event back.

---

### **Week 3: Event Streaming & Advanced Patterns**

_Goal: Move from simple queues to high-throughput event streaming and state management._

- **Day 15: Queues vs. Event Streams.**
  - _Study:_ Why Kafka is different from RabbitMQ. Transient messages vs. durable event logs.
- **Day 16: Apache Kafka Fundamentals.**
  - _Study:_ Topics, Partitions, Offsets, Brokers, and Consumer Groups.
  - _Action:_ Spin up Kafka via Docker. Run a basic producer and consumer CLI.
- **Day 17: Kafka in Practice.**
  - _Study:_ Managing consumer groups and scaling consumption.
  - _Action:_ Write a small app that publishes a stream of events (e.g., user clicks) to Kafka, and have multiple consumer instances read them.
- **Day 18: Redis Pub/Sub & In-Memory Messaging.**
  - _Study:_ When to use Redis for communication (fast, ephemeral) vs. Kafka/RabbitMQ.
- **Day 19: CQRS (Command Query Responsibility Segregation).**
  - _Study:_ Separating the "write" database from the "read" database and using events to keep them in sync.
  -
- **Day 20: Event Sourcing.**
  - _Study:_ Storing state as a sequence of events rather than just the current status (e.g., your bank account balance is a sum of all transactions, not just a single number).
- **Day 21: Week 3 Consolidation Project.**
  - _Action:_ Build a mini CQRS setup. `Service A` writes data to a Postgres DB and publishes an event to Kafka. `Service B` consumes the event and updates an Elasticsearch or MongoDB "read model."

---

### **Week 4: Resilience, Distributed Transactions & Security**

_Goal: Making sure your communication doesn't break when the real world gets messy._

- **Day 22: Distributed Transactions (Saga Pattern).**
  - _Study:_ You can't use standard database transactions across microservices. Learn the Saga Pattern (Choreography vs. Orchestration) to handle rollbacks/compensating transactions.
- **Day 23: The Transactional Outbox Pattern.**
  - _Study:_ How do you reliably save to your database _and_ publish an event to a queue without risking one failing? Look into the Outbox pattern and Change Data Capture (CDC / Debezium).
- **Day 24: Fault Tolerance (Circuit Breakers & Retries).**
  - _Study:_ Handling synchronous failures. The Circuit Breaker pattern (Closed, Open, Half-Open).
  - _Action:_ Implement a retry mechanism with exponential backoff in your code.
- **Day 25: Observability & Tracing.**
  - _Study:_ When a request hops through 5 services, how do you debug it? Learn about Correlation IDs and distributed tracing (Jaeger, OpenTelemetry).
- **Day 26: Service Mesh Overview.**
  - _Study:_ What are Istio and Linkerd? How they handle service-to-service communication, load balancing, and retries at the infrastructure layer (sidecar proxies).
- **Day 27: Security in Transit.**
  - _Study:_ Securing communication. mTLS (Mutual TLS) and passing JWTs (JSON Web Tokens) between services.
- **Day 28: Final Architecture Review.**
  - _Action:_ Draw a complete, production-ready architecture diagram for a hypothetical e-commerce store utilizing an API Gateway, synchronous reads, asynchronous checkout flows (Kafka/RabbitMQ), and an Outbox pattern.

---

### **A Quick Tip for Success**

Don't try to learn the deep, language-specific syntax of every tool right away. Focus heavily on **the concepts** (e.g., _why_ Kafka uses partitions, _why_ idempotency matters). The code can always be looked up later.
