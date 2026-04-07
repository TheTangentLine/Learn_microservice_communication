### **Day 7: Week 1 Consolidation Project**

Today, we are putting all the pieces together. You are going to build a fully functional, 3-tier Synchronous microservice architecture using Go and Docker Compose.

#### **The Architecture**

1. **User/Postman** sends an HTTP GET request to `http://localhost:8000/api/checkout?item=Nakroth`
2. **API Gateway** (Port 8000) receives it and reverse-proxies the HTTP request to the Order Service.
3. **Order Service** (Port 8081) receives the HTTP request. It then acts as a gRPC Client and makes a lightning-fast binary call to the Inventory Service.
4. **Inventory Service** (Port 50051) receives the gRPC call, checks its "database," and returns a boolean.
5. The response travels all the way back up the chain to the user.

#### **1. Project Setup**

Create a new folder: `week1-final`.
Copy over your `pb` folder and `inventory` folder from Day 5 exactly as they were.

Your structure should look like this:

```text
week1-final/
├── docker-compose.yml
├── pb/               # Your generated store.pb.go and store_grpc.pb.go
├── inventory/        # Day 5 gRPC Server code + Dockerfile
├── order/            # NEW: HTTP Server + gRPC Client + Dockerfile
└── gateway/          # NEW: API Gateway code + Dockerfile
```

#### **2. The Order Service (The Middleman)**

This service must now listen for HTTP traffic from the Gateway, but speak gRPC to the Inventory.
In `order/main.go`:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"week1-final/pb" // Update this to match your go.mod module name
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var inventoryClient pb.InventoryServiceClient

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	if item == "" {
		http.Error(w, "Missing 'item' parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Processing checkout for: %s", item)

	// Make the synchronous gRPC call to Inventory
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := inventoryClient.CheckStock(ctx, &pb.StockRequest{ItemName: item})
	if err != nil {
		http.Error(w, "Failed to check inventory", http.StatusServiceUnavailable)
		return
	}

	if resp.GetInStock() {
		w.Write([]byte(fmt.Sprintf("Success! %s is ordered.\n", item)))
	} else {
		w.Write([]byte(fmt.Sprintf("Sorry, %s is out of stock.\n", item)))
	}
}

func main() {
	// Connect to gRPC Inventory Service
	// Note: We use "inventory:50051" because Docker will resolve the container name
	conn, err := grpc.Dial("inventory:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial inventory service: %v", err)
	}
	defer conn.Close()

	inventoryClient = pb.NewInventoryServiceClient(conn)

	// Start HTTP Server for the Gateway to talk to
	http.HandleFunc("/api/checkout", checkoutHandler)
	log.Println("Order Service listening for HTTP on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
```

#### **3. The API Gateway**

We will build a simple Go Reverse Proxy. It intercepts traffic on port 8000 and forwards it to the Order service.
In `gateway/main.go`:

```go
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// The URL of our internal Order Service (using Docker's DNS)
	orderServiceURL, err := url.Parse("http://order:8081")
	if err != nil {
		log.Fatal(err)
	}

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(orderServiceURL)

	// Route all traffic hitting /api/checkout to the Order Service
	http.HandleFunc("/api/checkout", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Gateway routing request for: %s", r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	log.Println("API Gateway listening on public port :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
```

#### **4. The Docker Compose Glue**

Create a `Dockerfile` in `order/` and `gateway/` (just like the ones we used in Day 3). Then, in the root `docker-compose.yml`:

```yaml
version: "3.8"
services:
  inventory:
    build: ./inventory
    # We don't expose ports here! It is hidden completely from the host machine.

  order:
    build: ./order
    depends_on:
      - inventory
    # Hidden from host machine as well.

  gateway:
    build: ./gateway
    ports:
      - "8000:8000" # ONLY the Gateway is exposed to the outside world
    depends_on:
      - order
```

---

### **Actionable Task for Today**

1. Build the structure and write the code.
2. Run `docker-compose up --build`.
3. Open your browser and test: `http://localhost:8000/api/checkout?item=Nakroth%20Cybercore%20Skin` (Should succeed).
4. Test: `http://localhost:8000/api/checkout?item=Random%20Sword` (Should fail).

Notice how you **cannot** access the Order service or Inventory service directly from your browser anymore. The Gateway is the only open door.

---

### **End of Week 1 Review & Question**

Congratulations on finishing Week 1! You've built a real, functioning microservice cluster that uses both REST and gRPC, hidden securely behind a Gateway.

Take a breath, because tomorrow we start Week 2, and we are going to tear this synchronous architecture apart.

**Review Question to kick off Week 2:**
In our current Week 1 project, if 100,000 people try to buy the Nakroth skin at the exact same millisecond, the Gateway passes 100,000 HTTP requests to the Order Service, which opens 100,000 gRPC connections to the Inventory Service.
**Based on what we discussed on Day 2, what is going to happen to our system, and what specific tool will we introduce in Week 2 to fix it?**
