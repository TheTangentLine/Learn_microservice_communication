### **Day 6: API Gateways**

So far, we have our `Order Service` talking to our `Inventory Service`. But how does a web browser or a mobile app talk to the `Order Service`?

As we discussed earlier, browsers don't do native gRPC very well. Furthermore, you don't want to expose your internal microservices directly to the public internet.

#### **1. What is an API Gateway?**

An API Gateway is a server that acts as the single entry point into your system. It sits between your users and your microservices. It is essentially a "reverse proxy" on steroids.

#### **2. Why do we need it?**

- **Protocol Translation (The big one for us):** The Gateway can accept standard **HTTP/REST/JSON** requests from a web browser, and then translate them into lightning-fast **gRPC/Protobuf** calls to your internal services.
- **Routing:** It routes `api.myapp.com/orders` to the Order Service, and `api.myapp.com/users` to the User Service.
- **Authentication/Authorization:** Instead of putting JWT validation code in _every single_ microservice, you validate the user's token once at the Gateway. If the token is fake, the request is rejected before it ever touches your internal network.
- **Rate Limiting:** Stop DDoS attacks or abusive users at the front door.

#### **3. Popular API Gateways**

You rarely write an API Gateway from scratch. You use battle-tested open-source tools:

- **Kong:** Extremely popular, highly performant (built on NGINX/Lua).
- **Envoy:** The modern standard for cloud-native applications (often used as a sidecar, but works great as an edge gateway).
- **KrakenD:** Written in Go, incredibly fast, and very easy to configure for REST-to-REST or REST-to-gRPC.
- **Cloud Managed:** AWS API Gateway, Google Cloud API Gateway.

---

### **Actionable Task for Today**

Today is about design and prep for tomorrow's Week 1 final project.

Take out a piece of paper and draw the architecture we are going to build tomorrow:

1. Draw a box representing a **Client** (Browser/Postman). It sends an HTTP GET request to `/api/checkout`.
2. Draw a box representing the **API Gateway**. It receives the HTTP request.
3. Draw a box representing the **Order Service**. The Gateway routes the request here. (Decide now if you want the Gateway to talk to the Order Service via HTTP or gRPC).
4. Draw a box representing the **Inventory Service**. The Order Service talks to this via gRPC to check stock.

_Note: For tomorrow's code challenge, we won't use a heavy tool like Kong just yet. We will write a very simple custom API Gateway in Go or Python to deeply understand how the routing works, or we can use NGINX via Docker._

---

### **Day 6 Revision Question**

By putting an API Gateway in front of all our microservices, we have successfully hidden our internal network, centralized our authentication, and enabled protocol translation.

However, looking at the architecture you just drew, **what is the most obvious, glaring architectural risk we just introduced to our system by routing _all_ traffic through this single Gateway?** Let me know your thoughts!

**Answer:**
The absolute biggest risk: **Single Point of Failure (SPOF).**

If your Gateway goes down, it doesn't matter if you have 1,000 perfectly healthy microservices sitting behind it; your entire application is offline to the outside world.

_(How do we fix this in the real world? We run multiple instances of the API Gateway and put a highly available Cloud Load Balancer in front of them. If one Gateway container crashes, the Load Balancer just routes traffic to the other ones)._
