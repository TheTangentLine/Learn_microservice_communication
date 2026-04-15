### **Day 27: Security in Transit (mTLS & JWTs)**

Until now, our microservices have been talking over **plaintext**. Anyone on the internal network could intercept those packets and read credit card numbers. Today we lock the system down.

#### **1. Mutual TLS (mTLS)**

Standard TLS: your browser connects to your bank, the bank's server presents a certificate proving who it is, and traffic is encrypted. The bank doesn't ask _your browser_ for a certificate — just a password.

In **Mutual TLS**, both sides verify each other with cryptographic certificates:

```mermaid
sequenceDiagram
    participant OrderSvc as Order Service
    participant InventorySvc as Inventory Service

    OrderSvc->>InventorySvc: ClientHello + "I am Order Service"\n(presents certificate)
    InventorySvc->>InventorySvc: verify Order Service cert\nagainst trusted CA
    InventorySvc->>OrderSvc: ServerHello + "I am Inventory Service"\n(presents certificate)
    OrderSvc->>OrderSvc: verify Inventory Service cert\nagainst trusted CA
    OrderSvc->>InventorySvc: cryptographic handshake complete
    Note over OrderSvc,InventorySvc: All traffic from this point\nis encrypted end-to-end
```

**The Nightmare:** Manually rotating thousands of TLS certificates every 30 days across hundreds of services is a DevOps disaster.

**The Savior:** The Service Mesh. With Istio, the Envoy sidecars **automatically generate, rotate, and validate** mTLS certificates in the background. Your Go app still sends plaintext to `localhost` — Envoy handles all the encryption transparently.

#### **2. Passing Identity (JWT)**

mTLS proves _which microservice_ is calling. But how does the Inventory Service know _which user_ clicked the button?

We use **JSON Web Tokens (JWT)**.

```mermaid
sequenceDiagram
    participant User as User's Browser
    participant AuthSvc as Auth Service
    participant Gateway as API Gateway
    participant OrderSvc as Order Service
    participant InventorySvc as Inventory Service

    User->>AuthSvc: POST /login (username, password)
    AuthSvc->>AuthSvc: validate credentials
    AuthSvc-->>User: JWT { user_id: 99, role: "admin" }\n(signed with private key)

    User->>Gateway: POST /checkout\nAuthorization: Bearer <JWT>
    Gateway->>Gateway: verify JWT signature\n(no DB lookup needed)
    Gateway->>Gateway: check Redis blacklist for jti
    Gateway->>OrderSvc: forward request + JWT header
    OrderSvc->>InventorySvc: gRPC + pass JWT in metadata
    InventorySvc->>InventorySvc: decode JWT\n"User 99 is admin — allowed"
```

1. The user logs in. The Auth Service gives their browser a signed JWT: `{"user_id": 99, "role": "admin"}`.
2. The browser sends the JWT in the `Authorization` header to the API Gateway.
3. The Gateway validates the cryptographic signature. A fake token is rejected instantly — without a database query.
4. The validated JWT header is passed downstream to the Order Service, then to the Inventory Service.

Any microservice deep in the network can instantly decode the JWT and know the user's identity and role — without querying a central database.

---

### **Actionable Task for Today**

Open [jwt.io](https://jwt.io) in your browser and examine a JWT payload. Notice it is just Base64-encoded JSON — **anyone can decode and read it**. The security comes from the **Signature** at the bottom, which mathematically guarantees the data hasn't been tampered with.

In Go, look up the [`golang-jwt/jwt`](https://github.com/golang-jwt/jwt) library — the industry standard for parsing and validating JWTs in HTTP middleware.

---

### **Day 27 Revision Question**

JWTs are stateless — microservices validate them mathematically without a database. But imagine User 99 goes rogue and an admin clicks "BAN USER."

**If User 99's JWT doesn't expire for another 24 hours and your services validate it mathematically without checking the database, how do you stop them from using your system for the next 24 hours?**

**Answer: Two complementary strategies**

```mermaid
flowchart TB
    subgraph blacklist ["Strategy 1: Redis Blacklist at Gateway"]
        Ban1["Admin bans User 99"]
        Redis1[("Redis\nKey: jti_of_jwt_99\nTTL: matches JWT expiry)")]
        GW1["API Gateway\nO(1) Redis check before\npassing request internally"]
        Ban1 -->|"INSERT jti into Redis"| Redis1
        GW1 -->|"check"| Redis1
        Redis1 -->|"found = banned, reject"| GW1
    end

    subgraph shortlived ["Strategy 2: Short-Lived JWTs + Refresh Tokens"]
        JWT2["JWT expiry: 15 minutes"]
        Refresh["Refresh Token (stateful)\nsaved in DB"]
        AuthSvc2["Auth Service"]
        Ban2["Admin bans User 99\nin DB"]

        JWT2 -->|"expires"| AuthSvc2
        AuthSvc2 -->|"check Refresh Token"| Refresh
        Ban2 -->|"mark user banned"| Refresh
        Refresh -->|"banned = deny new JWT"| AuthSvc2
    end
```

1. **API Gateway Redis Blacklist:** When User 99 is banned, drop their JWT ID (`jti`) into Redis with a TTL matching the token's expiry. The Gateway does a blazing-fast O(1) Redis check before every request. Internal microservices remain completely stateless — they never check a database.

2. **Short-Lived JWTs + Refresh Tokens:** Keep JWT expiry very short (e.g., 15 minutes). A banned user can cause trouble for at most 14 minutes and 59 seconds. To avoid forcing users to re-login every 15 minutes, the frontend uses a stateful **Refresh Token** (saved in the database). When the JWT expires, the frontend quietly asks the Auth Service for a new one. If the user is banned, the Auth Service denies the refresh request — locking them out permanently.
