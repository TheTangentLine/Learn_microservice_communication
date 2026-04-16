# Learn Microservice Communication

A 28-day sprint covering how distributed services talk to each other — synchronously, asynchronously, and everything in between.

---

## Stack

- **Language:** Go
- **Brokers:** RabbitMQ, Apache Kafka, Redis Pub/Sub
- **Infrastructure:** Docker Compose, LocalStack (AWS SQS/SNS)
- **Patterns:** Outbox, Saga, CQRS, Event Sourcing, Circuit Breaker

---

## Week 1 — Synchronous Communication

_How services call each other directly and wait for a response._

| Day | Topic |
|-----|-------|
| [Day 1](Notes/Week1-Fundamentals_and_Synchronous_communication/day1-microservices_paradigm.md) | The Microservices Paradigm & The 8 Fallacies of Distributed Computing |
| [Day 2](Notes/Week1-Fundamentals_and_Synchronous_communication/day2-sync_vs_async.md) | Sync vs. Async — trade-offs and when to use each |
| [Day 3](Notes/Week1-Fundamentals_and_Synchronous_communication/day3-RESTful.md) | RESTful HTTP — building two services that talk over JSON |
| [Day 4](Notes/Week1-Fundamentals_and_Synchronous_communication/day4-RPC_and_gRPC.md) | RPC & gRPC — why binary Protobuf over HTTP/2 beats REST internally |
| [Day 5](Notes/Week1-Fundamentals_and_Synchronous_communication/day5-implementing_gRPC.md) | Implementing gRPC in Go — `.proto` contracts, timeouts, context propagation |
| [Day 6](Notes/Week1-Fundamentals_and_Synchronous_communication/day6-api_gateway.md) | API Gateways — the single entry point: routing, auth, rate limiting |
| [Day 7](Notes/Week1-Fundamentals_and_Synchronous_communication/day7-consolidation_project.md) | Project: Gateway → Order Service (HTTP) → Inventory Service (gRPC) |

---

## Week 2 — Asynchronous Communication & Message Queues

_How services communicate without waiting for each other._

| Day | Topic |
|-----|-------|
| [Day 8](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day8-EDA.md) | Event-Driven Architecture — commands vs. events, temporal decoupling |
| [Day 9](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day9-messages_queues_and_brokers.md) | Message Brokers — point-to-point queues vs. publish/subscribe |
| [Day 10](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day10-RabbitMQ_basics.md) | RabbitMQ Basics — AMQP, producers, consumers, manual ACKs |
| [Day 11](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day11-advanced_routing.md) | Exchanges & Routing — Fanout, Direct, Topic exchange types |
| [Day 12](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day12-cloud_native_queues.md) | Cloud Queues — AWS SQS, SNS, and the SNS-to-SQS fanout pattern |
| [Day 13](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day13-message_delivery_and_idempotency.md) | Delivery Guarantees — at-most-once, at-least-once, exactly-once, idempotency |
| [Day 14](Notes/Week2-Asynchronous_Communication_And_Message_Queues/day14-consolidation_project.md) | Project: Order Service publishes to RabbitMQ, Inventory Service consumes |

---

## Week 3 — Event Streaming & Advanced Patterns

_High-throughput streaming and separating reads from writes._

| Day | Topic |
|-----|-------|
| [Day 15](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day15-queues_vs_event_streams.md) | Queues vs. Streams — why Kafka's append-only log changes everything |
| [Day 16](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day16-Apache_Kafka_fundamentals.md) | Kafka Fundamentals — topics, partitions, offsets, consumer groups |
| [Day 17](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day17-Kafka_in_practice.md) | Kafka in Go — keyed messages, partition assignment, ordering guarantees |
| [Day 18](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day18-Redis_PubSub.md) | Redis Pub/Sub — ephemeral messaging and the WebSocket notification pattern |
| [Day 19](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day19-CQRS.md) | CQRS — separate write databases from read databases, synced by events |
| [Day 20](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day20-event_sourcing.md) | Event Sourcing — store events not state; replay to rebuild history |
| [Day 21](Notes/Week3-Event_Streaming_and_Advanced_Patterns/day21-consolidation_challenge.md) | Project: Command API writes to Kafka; Query Service builds a read model |

---

## Week 4 — Resilience, Transactions & Security

_Keeping communication reliable and safe when things go wrong._

| Day | Topic |
|-----|-------|
| [Day 22](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day22-Saga_pattern.md) | The Saga Pattern — distributed rollbacks via choreography or orchestration |
| [Day 23](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day23-Transactional_Outbox_pattern.md) | The Outbox Pattern — atomically save to DB and publish to a broker |
| [Day 24](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day24-Circuit_Breakers_and_retries.md) | Circuit Breakers & Retries — fail fast, recover gracefully |
| [Day 25](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day25-observability_and_distributed_tracing.md) | Distributed Tracing — follow one request across many services with Trace IDs |
| [Day 26](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day26-service_mesh_overview.md) | Service Mesh — move retries, circuit breaking, and tracing out of your code |
| [Day 27](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day27-mTLS_and_JWTs.md) | mTLS & JWTs — encrypt service-to-service traffic and propagate user identity |
| [Day 28](Notes/Week4-Resilience_and_Distributed_transactions_and_Security/day28-the_final_architecture_review.md) | Final Review — full end-to-end architecture for a high-load purchase flow |

---

## Extra — Error Handling & Observability

_How services communicate failures clearly — to clients, to other services, and to engineers debugging at 2am._

| | Topic |
|---|---|
| [Extra 1](Notes/Extra-Error_Handling_and_Observability/extra1-error_classification_and_propagation.md) | Error Classification — business errors vs infrastructure errors, and which layer owns each |
| [Extra 2](Notes/Extra-Error_Handling_and_Observability/extra2-structured_logging_with_trace_ids.md) | Structured Logging — log levels, Trace ID injection, what each layer logs and what it must not |
| [Extra 3](Notes/Extra-Error_Handling_and_Observability/extra3-go_error_patterns.md) | Go Error Patterns — typed sentinels, `%w` wrapping chains, three-layer handler, Trace middleware |

---

Focus on the **why** behind each pattern. The code is lookable — the mental models are what matter.
