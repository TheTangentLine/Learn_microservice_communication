### **Day 1: The Microservices Paradigm & The Hard Truths**

**1. The "Why": Monolith vs. Microservices**
Before services can communicate, you need to understand why they are split apart in the first place.

- **The Monolith:** All your business logic (users, inventory, billing) lives in one codebase and runs as a single process. Functions call other functions directly in memory. It's fast and easy to deploy initially, but becomes a nightmare to scale or update as teams grow.
- **Microservices:** The monolith is chopped up into independent, deployable services organized around business capabilities.
- **The Trade-off:** You are trading _codebase complexity_ for _infrastructure complexity_. Instead of a fast, guaranteed in-memory function call, Service A now has to reach out over a network to talk to Service B.

**2. The 8 Fallacies of Distributed Computing**
This is the most important concept to grasp today. When junior developers build microservices, they often assume the network behaves exactly like an internal computer component. It doesn't. In 1994, Peter Deutsch at Sun Microsystems drafted these 8 false assumptions programmers make when building distributed systems:

1.  **The network is reliable.** (Spoiler: Cables break, routers crash, AWS goes down).
2.  **Latency is zero.** (Sending data across the country takes time).
3.  **Bandwidth is infinite.** (You can't send a 5GB JSON payload instantly).
4.  **The network is secure.** (Assume traffic can be intercepted).
5.  **Topology doesn't change.** (Servers spin up and die constantly in Docker/Kubernetes).
6.  **There is one administrator.** (Different teams own different services).
7.  **Transport cost is zero.** (Data transfer costs money and compute power).
8.  **The network is homogeneous.** (Your Python service might be talking to a Go service on Linux, while receiving data from an iOS app).

Every pattern we learn over the next 4 weeks (queues, retries, circuit breakers) exists specifically to solve one of these 8 fallacies.

**3. Actionable Task: Environment Setup**
Let's get your machine ready for tomorrow.

1.  **Choose your weapon:** Install [Golang](https://go.dev/doc/install) or [Python](https://www.python.org/downloads/).
2.  **Install Docker:** Download and install [Docker Desktop](https://www.docker.com/products/docker-desktop/). We will rely heavily on `docker-compose` to spin up databases, API gateways, and message brokers easily.

---

### **Weekly Challenge Teaser**

At the end of Week 1 (Day 7), your challenge will be to build a locally running, 3-tier synchronous architecture using Docker Compose. You will have an API Gateway that routes HTTP traffic to an `Order Service` (written in Go or Python), which will synchronously call an `Inventory Service` to check stock before confirming the order.

---

### **Day 1 Revision Question**

Imagine you have an e-commerce Monolith that you just split into a `Checkout Service` and a `Payment Service`. They talk over HTTP. Based on the **8 Fallacies of Distributed Computing**, what are two specific things that could go wrong when the `Checkout Service` asks the `Payment Service` to process a credit card, which would never have happened in the old Monolith?

**Answer:**

1. **The network could be broken:** (Fallacy 1: The network is reliable). In a monolith, if the checkout code calls the payment code, it works unless the whole server is dead. In microservices, a router could glitch, a DNS lookup could fail, or the Payment Service container might be restarting, causing your HTTP request to drop into a black hole.
2. **Latency causing race conditions/timeouts:** (Fallacy 2: Latency is zero). If the Payment Service takes 5 seconds to process a card, your Checkout Service is left hanging for those 5 seconds. If thousands of users do this at once, your Checkout Service runs out of memory just waiting for replies.

You've got the mindset down. Let's move on to how we actually send those messages.
