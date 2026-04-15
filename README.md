# Learn Microservice Communication

A 28-day, hands-on sprint through the complete architecture of distributed systems — from synchronous HTTP all the way to enterprise-grade resilience patterns, event streaming, and security.

Every day includes: concept notes, a Mermaid architecture diagram, working Go code (where applicable), and a revision Q&A with full answers.

---

## The Stack

- **Language:** Go
- **Brokers:** RabbitMQ, Apache Kafka, Redis Pub/Sub
- **Infrastructure:** Docker Compose, LocalStack (AWS SQS/SNS)
- **Observability:** OpenTelemetry, Jaeger
- **Security:** mTLS (Istio/Envoy), JWT
- **Patterns:** Outbox, Saga, CQRS, Event Sourcing, Circuit Breaker

---

## Week 1 — Fundamentals & Synchronous Communication

_Goal: Understand why microservices exist, and how services talk to each other in real time._

| Day | Topic | Key Diagram |
|-----|-------|-------------|
| [Day 1](Notes/Week1-Fundamentals_and_Synchronous_communication/day1-microservices_paradigm.md) | The Microservices Paradigm & The 8 Fallacies | Monolith vs Microservices flowchart |
| [Day 2](Notes/Week1-Fundamentals_and_Synchronous_communication/day2-sync_vs_async.md) | Sync vs. Async Overview | Sequence diagram: phone call vs email |
| [Day 3](Notes/Week1-Fundamentals_and_Synchronous_communication/day3-RESTful.md) | RESTful Communication (Go code) | HTTP call sequence with timeout paths |
| [Day 4](Notes/Week1-Fundamentals_and_Synchronous_communication/day4-RPC_and_gRPC.md) | RPC & gRPC Fundamentals | REST vs gRPC comparison + proto codegen |
| [Day 5](Notes/Week1-Fundamentals_and_Synchronous_communication/day5-implementing_gRPC.md) | Implementing gRPC in Go | gRPC call sequence + context/timeout |
| [Day 6](Notes/Week1-Fundamentals_and_Synchronous_communication/day6-api_gateway.md) | API Gateways | Gateway routing + Load Balancer SPOF |
| [Day 7](Notes/Week1-Fundamentals_and_Synchronous_communication/day7-consolidation_project.md) | Week 1 Project: 3-Tier Sync Architecture | Full architecture: Gateway → Order → Inventory |

**Week 1 Project:** Gateway (port 8000) routes HTTP to an Order Service, which calls an Inventory Service via gRPC. Only the Gateway is exposed — internal services are hidden inside the Docker network.

---

## Week 2 — Asynchronous Communication & Message Queues

_Goal: Decouple services so they don't have to wait for each other._

| Day | Topic | Key Diagram |
|-----|-------|-------------|
| [Day 8](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day8-EDA.md) | Intro to Event-Driven Architecture | EDA Pub/Sub flowchart + Commands vs Events |
| [Day 9](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day9-messages_queues_and_brokers.md) | Message Queues & Brokers | Point-to-Point vs Pub/Sub |
| [Day 10](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day10-RabbitMQ_basics.md) | RabbitMQ Basics (Go code) | AMQP sequence + auto-ack danger |
| [Day 11](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day11-advanced_routing.md) | Advanced Routing: Exchanges | Fanout / Direct / Topic exchange types |
| [Day 12](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day12-cloud_native_queues.md) | Cloud-Native Queues (AWS SQS/SNS) | SNS-to-SQS fanout + Visibility Timeout |
| [Day 13](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day13-message_delivery_and_idempotency.md) | Message Delivery Guarantees & Idempotency | State machine: PENDING → COMPLETED |
| [Day 14](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day14-consolidation_project.md) | Week 2 Project: Async EDA Architecture | Full async system with manual ACKs |

**Week 2 Project:** Order Service responds HTTP 202 immediately, publishes `OrderPlaced` to a RabbitMQ Fanout Exchange. Inventory Service consumes with manual ACKs and idempotency key enforcement.

---

## Week 3 — Event Streaming & Advanced Patterns

_Goal: High-throughput streaming and advanced data architecture patterns._

| Day | Topic | Key Diagram |
|-----|-------|-------------|
| [Day 15](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day15-queues_vs_event_streams.md) | Queues vs. Event Streams | Queue (delete on read) vs Stream (consumer owns offset) |
| [Day 16](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day16-Apache_Kafka_fundamentals.md) | Apache Kafka Fundamentals | Topics / Partitions / Offsets / Consumer Groups |
| [Day 17](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day17-Kafka_in_practice.md) | Kafka in Practice (Go code) | Hash-keyed producer + idle worker scenario |
| [Day 18](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day18-Redis_PubSub.md) | Redis Pub/Sub & In-Memory Messaging | Relay Pattern: Kafka → DB + Redis → WebSocket |
| [Day 19](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day19-CQRS.md) | CQRS | Command/Query separation flowchart |
| [Day 20](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day20-event_sourcing.md) | Event Sourcing | Event replay + Crypto-Shredding (GDPR) |
| [Day 21](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day21-consolidation_challenge.md) | Week 3 Project: CQRS + Event Sourcing | Command API → Kafka → Query Service read model |

**Week 3 Project:** Command Service writes price-change events to Kafka (Event Sourcing). Query Service consumes those events and builds an in-memory read model (CQRS). No shared database between them.

---

## Week 4 — Resilience, Distributed Transactions & Security

_Goal: Make the system survive the real world — failures, fraud, and distributed rollbacks._

| Day | Topic | Key Diagram |
|-----|-------|-------------|
| [Day 22](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day22-Saga_pattern.md) | The Saga Pattern | Choreography vs Orchestration + Compensating Tx |
| [Day 23](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day23-Transactional_Outbox_pattern.md) | The Transactional Outbox Pattern | Outbox ACID write + Go worker pool |
| [Day 24](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day24-Circuit_Breakers_and_retries.md) | Circuit Breakers & Retries | State diagram: CLOSED/OPEN/HALF-OPEN |
| [Day 25](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day25-observability_and_distributed_tracing.md) | Observability & Distributed Tracing | Trace ID sequence + Span Gantt waterfall |
| [Day 26](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day26-service_mesh_overview.md) | Service Mesh (Istio/Envoy) | Before/after Sidecar + Kubernetes Pod |
| [Day 27](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day27-mTLS_and_JWTs.md) | mTLS & JWTs | mTLS handshake + JWT lifecycle + revocation |
| [Day 28](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day28-the_final_architecture_review.md) | Final Architecture Review | Full system + 7-step purchase sequence |

**Week 4 Final:** End-to-end architecture handling 500,000 concurrent users for a game skin launch — JWT auth, Outbox pattern, Kafka fan-out, Circuit Breaker against Stripe, Idempotent inventory unlock, Redis Pub/Sub notification, and a real-time CQRS analytics dashboard.

---

## A Note on Learning

Don't try to memorize syntax for every tool. Focus on **why** each pattern exists — why Kafka uses partitions, why idempotency is non-negotiable, why the Circuit Breaker has three states. The code is always lookable. The mental models are what make you a senior engineer.
