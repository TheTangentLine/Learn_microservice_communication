### **Day 3: RESTful Communication (Writing the Code)**

Today, we are going to build a Synchronous HTTP interaction. We will create two microservices using **Golang** and connect them using Docker Compose.

Service A (`Order Service`) will receive a request from a user, and it will synchronously call Service B (`Inventory Service`) via standard HTTP/REST to ask, "Do we have this in stock?"

#### **1. The Project Structure**

Create a new folder on your computer called `day3-sync` and set up this exact structure:

```text
day3-sync/
├── docker-compose.yml
├── inventory/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
└── order/
    ├── Dockerfile
    ├── go.mod
    └── main.go
```

#### **2. The Inventory Service (Service B)**

This service simply listens for a request and replies with JSON.
In `inventory/go.mod`, put:

```go
module inventory
go 1.21
```

In `inventory/main.go`, put:

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func checkStockHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received stock check request...")
	w.Header().Set("Content-Type", "application/json")
	// Hardcoding true for today's example
	json.NewEncoder(w).Encode(map[string]bool{"in_stock": true})
}

func main() {
	http.HandleFunc("/check-stock", checkStockHandler)
	log.Println("Inventory Service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
```

In `inventory/Dockerfile`, put:

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["./main"]
```

#### **3. The Order Service (Service A)**

This service receives a user request and _makes_ an HTTP call to the Inventory Service.
In `order/go.mod`, put:

```go
module order
go 1.21
```

In `order/main.go`, put:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type StockResponse struct {
	InStock bool `json:"in_stock"`
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting order creation process...")

	// 1. Make a Synchronous HTTP GET request to the Inventory Service
	// Notice the URL uses "inventory", which Docker will resolve to the other container!
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://inventory:8081/check-stock")

	if err != nil {
		http.Error(w, "Failed to reach Inventory Service: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// 2. Parse the JSON response
	var stockResp StockResponse
	json.NewDecoder(resp.Body).Decode(&stockResp)

	// 3. Respond to the user
	if stockResp.InStock {
		w.Write([]byte("Success! Order placed for 1x Millennium Falcon Lego Set.\n"))
	} else {
		w.Write([]byte("Failed: Item out of stock.\n"))
	}
}

func main() {
	http.HandleFunc("/create-order", createOrderHandler)
	log.Println("Order Service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

In `order/Dockerfile`, put:

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["./main"]
```

#### **4. The Docker Compose Glue**

This file spins up both Go applications and puts them on the same virtual network so they can talk to each other.
In the root `docker-compose.yml`, put:

```yaml
version: "3.8"
services:
  inventory:
    build: ./inventory
    ports:
      - "8081:8081"

  order:
    build: ./order
    ports:
      - "8080:8080"
    depends_on:
      - inventory
```

---

### **Actionable Task for Today**

1. Open your terminal in the `day3-sync` folder.
2. Run: `docker-compose up --build`
3. Open a new terminal tab or your web browser and go to: `http://localhost:8080/create-order`
4. Look at the logs in your Docker terminal. You will see the Order service receive your request, reach out over the Docker network to the Inventory service, get the JSON response, and send the final success message back to you.

---

### **Day 3 Revision Question**

We just built a synchronous system. While your `docker-compose` is running, open a new terminal and run `docker-compose pause inventory`. This simulates the Inventory service crashing or freezing.

Now, try to hit `http://localhost:8080/create-order` again. **What happens, and how does this perfectly demonstrate the main flaw of Synchronous communication we discussed on Day 2?**

**Output:** `Failed to reach Inventory Service: Get "http://inventory:8081/check-stock": context deadline exceeded (Client.Timeout exceeded while awaiting headers)`

**Explanation:** Because the `Order Service` is actively waiting for a response, if the `Inventory Service` is down (or paused), the whole operation fails and returns "Service Unavailable." This is the ultimate flaw of synchronous communication: **Temporal Coupling**. Both services must be alive, healthy, and fast at the exact same moment for the system to work.

If your system has 5 microservices that call each other synchronously, and each one has a 99% uptime guarantee, your overall system uptime drops to about 95% (0.99 ^ 5).

But for things like checking inventory during checkout, we often _have_ to be synchronous. So, if we are forced to be synchronous, how do we make it faster and more efficient than standard HTTP/REST?
