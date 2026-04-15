### **Day 19: CQRS (Command Query Responsibility Segregation)**

By combining Kafka, databases, and Redis, we just accidentally designed a CQRS system.

#### **1. What is CQRS?**

In traditional apps, you use the exact same database (and often the exact same code models) to _write_ data and to _read_ data.
CQRS says: **"Writing data is fundamentally different from Reading data. They should be split into two entirely different systems."**

- **The Command Side (Writes):** Focuses on business rules, validation, and safely storing the truth. It is usually a relational database (Postgres/MySQL) or an event stream (Kafka).
- **The Query Side (Reads):** Focuses entirely on speed and UI presentation. It is usually a NoSQL database (MongoDB), a search engine (Elasticsearch), or an in-memory cache (Redis).

#### **2. How do they stay in sync?**

Events!

1. A user places an order (Command). It goes into Postgres.
2. The Order Service publishes an `OrderUpdated` event to Kafka.
3. A separate worker reads Kafka and updates an Elasticsearch document (Query).
4. When the user visits their dashboard, the UI queries Elasticsearch (which is lightning fast and handles text searching beautifully), totally bypassing the heavy Postgres database.

---

### **Actionable Task for Today**

Grab your notepad and map out a CQRS architecture for an **E-commerce Product Search Page**.

1.  **The Command Side:** An Admin updates the price of the "Nakroth Skin" in the core inventory SQL database.
2.  **The Event:** How does that price change get into the broker?
3.  **The Query Side:** Your frontend users don't query the SQL database to search for items. They query a highly optimized **Elasticsearch** (or MongoDB) index. How does the event get from the broker into Elasticsearch?

---

### **Day 19 Revision Question**

CQRS gives us amazing read speeds and scaling, but it introduces our old friend: **Eventual Consistency**.

Imagine an Admin goes into the Command system and changes the price of an item from $10 to $15.
They click "Save", the UI says "Saved!", and they immediately refresh the public product page.

Because the event has to travel through Kafka and update the Query database, the public page still says "$10" for a few seconds.
**If you are the lead engineer, how do you handle this UI/UX problem where the user sees stale data immediately after making a write?**

**Answer:**

As a backend engineer, it is sometimes hard to admit, but **the best solution to an eventual consistency problem is often a UI trick.**

If you can't bend the laws of physics to make the network instantly fast, you manage the user's expectations instead.

There are two primary ways frontend teams handle this "Pending" state in a CQRS architecture:

1. **Honest UI:** The Admin clicks "Save." The UI immediately shows a toast notification saying _"Price update queued..."_ or puts a little spinning sync icon next to the price. It stays that way until the UI receives a WebSocket ping confirming the Elasticsearch database has been updated.
2. **Optimistic UI:** The UI _assumes_ the backend will succeed. It instantly changes the text on the screen to `$15` so it feels blazing fast to the user. Behind the scenes, the UI waits for the backend confirmation. If Kafka crashes and the backend returns a failure 5 seconds later, the UI reverts the price back to `$10` and pops up an error: _"Failed to save, please try again."_
