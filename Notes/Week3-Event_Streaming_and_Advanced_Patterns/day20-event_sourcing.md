### **Day 20: Event Sourcing**

We talked about Event Sourcing briefly at the end of Week 2. Today, we define it properly. It is the ultimate evolution of Event-Driven Architecture.

#### **1. The Flaw of Traditional Databases**

In traditional systems (CRUD: Create, Read, Update, Delete), your database only stores the **current state**.
If User A's bank balance is `$500`, the database literally says `balance = 500`.
If they withdraw `$100`, you run an `UPDATE` command, overwriting the old data so it now says `balance = 400`.

**The Problem:** You just permanently destroyed the history. If you want to know _why_ the balance is `$400`, the database can't tell you. (You usually have to build messy, separate "audit log" tables to track this).

#### **2. The Event Sourcing Paradigm**

Event Sourcing says: **Stop storing the current state. Only store the events that happened.**

Think of a bank ledger:

1. `AccountCreated` (+$0)
2. `Deposited` (+$500)
3. `Withdrew` (-$100)

If you want to know the user's current balance, you don't query a `balance` column. You write a function that grabs all their events and replays them from start to finish (`0 + 500 - 100 = 400`).

#### **3. Why is this so powerful?**

- **100% Auditability:** You have a mathematically perfect history of everything that ever happened in your system. You can never lose data to a bad `UPDATE` query.
- **Time Travel:** You can rebuild the exact state of the system at any point in time. "What was this user's cart like last Tuesday at 4:00 PM?" Just replay the events up to that exact timestamp and stop!
- **Perfect for Kafka & CQRS:** Kafka is the ultimate append-only event log. You use Kafka as your source of truth (Event Sourcing), and you replay those events to build fast read-models in Elasticsearch or Redis (CQRS).

#### **4. The Catch (Snapshots)**

Replaying 10,000 events every time a user logs in to check their bank balance is terribly slow.
To fix this, systems use **Snapshots**. Every 100 events, you save a "Snapshot" of the current state (`balance = 400`). Next time you need the balance, you load the Snapshot, and only replay the events that happened _after_ it.

---

### **Actionable Task for Today**

We are going to do a quick mental mapping exercise to prepare for tomorrow's Week 3 project.

Imagine a **Shopping Cart Service** built with Event Sourcing.

1. Write down the sequence of at least 4 specific Events that would happen if a user added a Nakroth Skin to their cart, added a Health Potion, changed their mind and removed the Potion, and then checked out.
2. What does the "Current State" look like after all 4 events are replayed?

---

### **Day 20 Revision Question**

Event Sourcing relies on an "Append-Only" database (like Kafka or a special Event Store DB). You can literally never run an `UPDATE` or `DELETE` command on past events.

However, in the real world, modern privacy laws like GDPR or the CCPA enforce the **"Right to be Forgotten."** A user has the legal right to demand you permanently delete all their personal data (Name, Email, Address) from your systems.

**If your core architecture relies on an immutable, un-deletable log of events from the past 5 years... how do you legally comply with a user's request to delete their personal data without destroying your Event Sourced history?**

**Answer:**

Instinct to use a "blacklist" is actually a fantastic first step, and it is a pattern companies use! You maintain a table of `deleted_user_ids`, and whenever your system replays events, it checks that list. If the event belongs to a blacklisted user, the system just skips it.

Meanwhile, the "sliding window" hits a legal wall. Under GDPR, you usually have 30 days to delete the data. Because an Event Sourced database is meant to be the permanent history of the company, its sliding window is usually set to **Infinite**. If you just hide the data but keep it on your hard drive forever, European regulators will still fine you heavily!

So, how do we delete data from an un-deletable database? The industry uses two brilliant workarounds.

### **1. Crypto-Shredding (The Most Elegant Solution)**

Instead of trying to delete the event, you make the event completely unreadable.

- When User 123 registers, you generate a unique cryptographic encryption key just for them. You store this key in a secure Key Management Database.
- Whenever you publish an event to Kafka (`OrderPlaced`), you **encrypt** the Personal Identifiable Information (PII) like their name and address using their unique key. The event just looks like gibberish in the Kafka log.
- **To delete the user:** You don't touch Kafka at all. You simply delete their encryption key from the Key Management Database. Instantly, all their historical events scattered across 5 years of logs become permanently, mathematically unrecoverable.

### **2. Log Compaction & Tombstones (The Kafka Way)**

If you are using Kafka, you can use the Log Compaction feature we briefly touched on during Day 15.

- You publish a new event to Kafka with the user's ID as the Key, but you set the message payload to `null`.
- This is called a **Tombstone** message.
- When Kafka runs its background cleanup process, it sees the Tombstone and says, _"Ah, the final state for this key is null!"_ and it will actively hunt down and delete all previous events with that key to save space.
