### **Day 12: Cloud-Native Queues (AWS SQS/SNS)**

Running RabbitMQ yourself is powerful, but it comes with a cost: you have to maintain the server, patch the OS, monitor the disk space, and handle scaling if you suddenly get a million users.

Most modern startups and enterprise teams prefer to use **Managed Cloud Services**. Today, we look at the undisputed king of cloud messaging: AWS.

#### **1. SQS (Simple Queue Service)**

SQS is Amazon's fully managed message queue.

- **The Good:** You don't manage any servers. It scales infinitely and automatically. You just pay a few cents for every million messages you send.
- **The Bad:** It is strictly **Point-to-Point**. It does _not_ have the routing power or "Exchanges" that RabbitMQ has. If you put a message in an SQS queue, one worker will read it, and it's gone.

#### **2. SNS (Simple Notification Service)**

Because SQS can't do Pub/Sub (Publish/Subscribe), AWS created SNS.

- SNS is essentially a **Fanout Exchange**.
- You publish a message to an SNS "Topic".

#### **3. The "SNS-to-SQS Fanout" Pattern**

This is one of the most famous architectural patterns in the cloud. Because SQS is great at holding messages safely, and SNS is great at copying messages to multiple places, we combine them to recreate what we built in Day 11.

1. Your Order Service publishes an `OrderPlaced` event to an **SNS Topic**.
2. You have two **SQS Queues**: an `Inventory Queue` and an `Email Queue`.
3. You "subscribe" both SQS Queues to the SNS Topic.
4. When the message hits SNS, it instantly drops a copy of the message into both SQS queues.
5. Your Go workers independently read from their respective SQS queues.

---

### **Actionable Task for Today**

You don't need a real AWS account or a credit card for this! We are going to use a brilliant tool called **LocalStack**, which mimics AWS locally inside Docker.

**1. Update your `docker-compose.yml`:**
Create a new folder `day12-aws`, and add this `docker-compose.yml`:

```yaml
version: "3.8"
services:
  localstack:
    image: localstack/localstack
    ports:
      - "4566:4566" # The main port for all fake AWS services
    environment:
      - SERVICES=sqs,sns
      - DOCKER_HOST=unix:///var/run/docker.sock
```

Run `docker-compose up -d`.

**2. Install the AWS CLI (Command Line Interface):**
Download the AWS CLI for your OS. Once installed, configure it with dummy credentials (since we are using localstack):

```bash
aws configure
# AWS Access Key ID: test
# AWS Secret Access Key: test
# Default region name: us-east-1
# Default output format: json
```

**3. Build the Architecture via Terminal:**
Let's build the SNS-to-SQS pattern! Run these commands one by one to see how cloud infrastructure is provisioned:

_Create the Queues:_

```bash
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name inventory-queue
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name email-queue
```

_Create the SNS Topic:_

```bash
aws --endpoint-url=http://localhost:4566 sns create-topic --name order-topic
```

_Subscribe the Queues to the Topic:_
_(Note: You'll need the ARNs—Amazon Resource Names—from the output of the previous commands, but they usually look like `arn:aws:sqs:us-east-1:000000000000:inventory-queue`)_

```bash
aws --endpoint-url=http://localhost:4566 sns subscribe \
    --topic-arn arn:aws:sns:us-east-1:000000000000:order-topic \
    --protocol sqs \
    --notification-endpoint arn:aws:sqs:us-east-1:000000000000:inventory-queue
```

_(You can skip writing the Go code today, just grasp how the cloud components link together!)_

---

### **Day 12 Revision Question**

AWS SQS has a core feature called the **Visibility Timeout**.

When your Payment Worker pulls a message from SQS, SQS doesn't delete it immediately. Instead, it makes the message "invisible" to all other workers for a default of 30 seconds. If the worker doesn't explicitly delete the message within those 30 seconds, SQS makes it visible again for someone else to pick up.

**Why is this concept identical to something we learned in RabbitMQ on Day 10, and what terrible thing happens to your system if your Go Payment Worker actually takes 45 seconds to process a credit card?**

**Answer:**

**1. What is it identical to in RabbitMQ?**
SQS's "Visibility Timeout" is identical to RabbitMQ's **Manual Acknowledgments (Ack/Nack)**. In RabbitMQ, the message stays in the queue (unacknowledged) until the worker explicitly says "I'm done" (`Ack`). If the worker dies or disconnects before acking, RabbitMQ requeues it.

**2. What terrible thing happens if processing takes 45 seconds?**
If your Visibility Timeout is 30 seconds, but your Payment Worker takes 45 seconds, the message becomes "visible" again at the 30-second mark.

- A second Payment Worker will see the message and start processing it.
- **The terrible result:** You process the same credit card twice, double-charging your customer!

To fix this, you either need to increase the default Visibility Timeout for that specific queue to be safely longer than your maximum expected processing time, or have your worker periodically call AWS to extend the timeout while it's still working.
