### **Day 24: Fault Tolerance (Circuit Breakers & Retries)**

For the past two weeks we solved problems with queues and events. But in the real world, you can't always be asynchronous. Sometimes your `Checkout Service` absolutely _must_ make a synchronous HTTP call to Stripe to charge a credit card.

What happens if Stripe's API goes down?

#### **1. The Problem: Cascading Failures**

If Stripe is down and your `Checkout Service` waits 10 seconds per timeout, 10,000 concurrent users will spawn 10,000 Goroutines — all waiting 10 seconds. Your server runs out of memory and crashes. Even worse, when Stripe comes back online, your servers hammer it with 10,000 retries simultaneously — crashing it again. This is a **Retry Storm**.

#### **2. The Solution: The Circuit Breaker Pattern**

Borrowed from electrical engineering — a proxy wrapper around your HTTP calls that protects the system.

```mermaid
stateDiagram-v2
    [*] --> CLOSED

    CLOSED --> CLOSED : request succeeds
    CLOSED --> OPEN : failure rate exceeds threshold\n(e.g. 50% fail in 10s)

    OPEN --> OPEN : all requests instantly rejected\n"Stripe is down" — no network call made
    OPEN --> HALFOPEN : cooldown period expires\n(e.g. 30 seconds)

    HALFOPEN --> CLOSED : probe request succeeds\nbreaker resets
    HALFOPEN --> OPEN : probe request fails\ncooldown restarts
```

- **CLOSED (Green Light):** Everything is healthy. Requests flow through normally. The breaker counts failures.
- **OPEN (Red Light):** Failure threshold crossed. The breaker instantly blocks all new requests and returns an error — no network call is made. This saves CPU/Memory and gives Stripe time to recover.
- **HALF-OPEN (Yellow Light):** After the cooldown, exactly _one_ probe request is let through. If it succeeds → CLOSED. If it fails → OPEN again.

#### **3. Adding Retries with Exponential Backoff + Jitter**

When the circuit is CLOSED, transient network blips still happen. Retry failed requests — but do it politely.

```mermaid
flowchart LR
    Req["Request fails"]
    R1["Retry in ~1.2s"]
    R2["Retry in ~2.5s"]
    R3["Retry in ~4.8s"]
    R4["Retry in ~9.3s"]
    GiveUp["Give up"]
    Success["Success"]

    Req --> R1
    R1 -->|"fails"| R2
    R2 -->|"fails"| R3
    R3 -->|"fails"| R4
    R4 -->|"fails"| GiveUp
    R1 -->|"succeeds"| Success
    R2 -->|"succeeds"| Success
    R3 -->|"succeeds"| Success
```

- **Bad:** Retry instantly → fail → retry instantly → fail. (DDoS your own dependency.)
- **Good (Exponential Backoff):** Retry in 1s → 2s → 4s → 8s → give up. Each delay doubles.
- **Pro-tip — Jitter:** Always add randomness. If 10,000 servers all retry at exactly 2.0s, they will collectively DDoS the target. Adding jitter (e.g., retry in 1.2s, then 2.5s) spreads the load.

---

### **Actionable Task for Today**

In Go, the industry-standard library is **`sony/gobreaker`**. Read the GitHub README for [`sony/gobreaker`](https://github.com/sony/gobreaker). Look at how you define the rules (`MaxRequests`, `Interval`, `Timeout`) and then wrap your `http.Get()` call inside an `Execute()` block.

---

### **Day 24 Revision Question**

When a Circuit Breaker trips to OPEN, your Checkout Service instantly fails fast — it refuses to call the `Fraud Detection API`. You now have a choice.

**Should you return a hard HTTP 500 Error to the user, or is there a better pattern you could implement when the breaker is open?**

**Answer: Graceful Degradation (Fallbacks)**

```mermaid
flowchart TD
    Request["Checkout Request"]
    CB{"Circuit Breaker\nState?"}
    Closed["CLOSED — call Fraud API normally"]
    Open["OPEN — breaker tripped"]

    FallbackQueue["Fallback A: Queue for later\nReturn 202 Pending\nAsync fraud check when API recovers"]
    FallbackRisk["Fallback B: Static Risk Rules\nOrder < $50 → auto-approve\nOrder > $500 → queue or reject"]

    Request --> CB
    CB -->|"healthy"| Closed
    CB -->|"Fraud API is down"| Open
    Open --> FallbackQueue
    Open --> FallbackRisk
```

Instead of failing the user's checkout entirely, you **gracefully degrade**:

1. **Async Queue Fallback:** Tell the frontend "Pending." Place the order in a queue. When the Fraud API recovers, a background worker drains the queue, runs the checks, and releases the orders.

2. **Static Risk Threshold Fallback:** If the Fraud API is down, the code says: _"If the order is under $50, auto-approve it and take the risk. If it's over $500, queue it or reject it."_

Both approaches serve the user a meaningful response instead of a cryptic 500 error — and they protect your system from a thundering herd when the Fraud API comes back online.
