### **Day 5: Implementing gRPC in Go**

Today, we are replacing the HTTP/REST communication from Day 3 with our high-speed gRPC setup.

Create a new folder called `day5-grpc`. Inside, we need a shared `pb` (protobuf) folder, an `inventory` folder, and an `order` folder.

#### **1. Generating the Go Code from Protobuf**

In your `day5-grpc` folder, create a file named `store.proto`:

```protobuf
syntax = "proto3";

// This tells protoc to put the generated code in a folder named "pb"
option go_package = "./pb";

package store;

message StockRequest {
  string item_name = 1;
}

message StockResponse {
  bool in_stock = 1;
}

service InventoryService {
  rpc CheckStock (StockRequest) returns (StockResponse);
}
```

Now, open your terminal in the `day5-grpc` folder and run this exact command to generate your Go code:

```bash
protoc --go_out=. --go-grpc_out=. store.proto
```

If successful, you will see a new `pb` folder appear with two files: `store.pb.go` and `store_grpc.pb.go`. **Never edit these files manually.** They are the strict contract both your services will use.

#### **2. The Inventory Service (gRPC Server)**

This service will implement the interface we just generated.
In `inventory/main.go`:

```go
package main

import (
	"context"
	"log"
	"net"

	// Replace "day5-grpc" with your actual go module name if different
	"day5-grpc/pb"
	"google.golang.org/grpc"
)

// server is used to implement pb.InventoryServiceServer
type server struct {
	pb.UnimplementedInventoryServiceServer
}

// CheckStock implements the RPC method defined in our .proto file
func (s *server) CheckStock(ctx context.Context, req *pb.StockRequest) (*pb.StockResponse, error) {
	log.Printf("Received check for: %s", req.GetItemName())

	// Let's pretend we are checking a database here.
	inStock := false
	if req.GetItemName() == "Nakroth Cybercore Skin" {
		inStock = true
	}

	return &pb.StockResponse{InStock: inStock}, nil
}

func main() {
	// 1. Open a TCP listener on port 50051 (standard gRPC port)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. Create a new gRPC server
	s := grpc.NewServer()

	// 3. Register our specific Inventory Service with the gRPC server
	pb.RegisterInventoryServiceServer(s, &server{})
	log.Printf("gRPC Inventory Server listening at %v", lis.Addr())

	// 4. Start serving
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

#### **3. The Order Service (gRPC Client)**

This service will "dial" the gRPC server and call the function as if it were local code.
In `order/main.go`:

```go
package main

import (
	"context"
	"log"
	"time"

	"day5-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Set up a connection to the server.
	// We use insecure credentials here because we aren't setting up SSL/TLS certificates today.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 2. Create the client stub using the generated code
	client := pb.NewInventoryServiceClient(conn)

	// 3. Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	item := "Nakroth Cybercore Skin"
	log.Printf("Attempting to order: %s", item)

	// THIS IS THE MAGIC! It looks like a local function call, but it's a network request.
	r, err := client.CheckStock(ctx, &pb.StockRequest{ItemName: item})
	if err != nil {
		log.Fatalf("could not check stock: %v", err)
	}

	if r.GetInStock() {
		log.Println("Success! Item is in stock. Proceeding with order.")
	} else {
		log.Println("Failed. Item is out of stock.")
	}
}
```

---

### **Actionable Task for Today**

1. Ensure your Go module is initialized (`go mod init day5-grpc` in the root folder) and run `go mod tidy` to download the gRPC dependencies.
2. Open two separate terminal windows.
3. In Terminal 1, run the server: `go run inventory/main.go`
4. In Terminal 2, run the client: `go run order/main.go`
5. Watch the client successfully trigger the Go code running in the server over a binary TCP connection! Change the `item` string in the `order/main.go` file to something else, re-run the client, and watch it fail the stock check.

---

### **Day 5 Revision Question**

Look at this line in the Order Service (Client) code:
`ctx, cancel := context.WithTimeout(context.Background(), time.Second)`

We are passing a `context` with a 1-second timeout to our gRPC `CheckStock` call. If the Inventory Service takes 3 seconds to respond (maybe the database is slow), what exactly happens to the Order Service, and **why is it critical to always include timeouts in synchronous microservice communication?** Let me know how the code runs on your machine and what your thoughts are on the revision question!

**Answer:**
The Order Service will abort the request after 1 second and throw an error. Meanwhile, the Inventory Service might keep processing the request in the background (wasting CPU and database connections), but the Order Service has already moved on. This is exactly how we prevent **cascading failures**—if one service is slow, it doesn't drag the rest of the system down with it.

### **How Go Context Translates to Python in gRPC**

Go's `context.Context` is a specific Go concept, and you cannot send a Go object over the network to a Python server.

Here is what actually happens under the hood:

1. **The Translation:** When your Go client calls `CheckStock` with a 1-second timeout context, the gRPC library intercepts it. It doesn't send the "context object." Instead, it looks at the deadline, calculates the time remaining, and attaches a special **HTTP/2 Header** (metadata) called `grpc-timeout`.
2. **Over the Wire:** The binary data and that `grpc-timeout` header fly over the network to the Python server.
3. **The Python Side:** The Python gRPC library receives the HTTP/2 request. It reads the `grpc-timeout` header. Python doesn't have Go's `context`, but the Python gRPC library automatically generates a **`grpc.ServicerContext`** object for you and passes it into your Python function.

If you wrote the Inventory Server in Python, the function signature would look like this:

```python
def CheckStock(self, request, context):
    # 'context' here is a Python grpc.ServicerContext object

    # You can check how much time the Go client gave you!
    time_left = context.time_remaining()

    # You can check if the Go client already cancelled/timed out
    if not context.is_active():
        return # Stop processing, the client gave up!
```

So, gRPC perfectly translates the _meaning_ of the context (the timeout, cancellation signals, and metadata) between completely different programming languages using HTTP/2 headers as the middleman.

Ready for the next step? Let's give our system a front door.
