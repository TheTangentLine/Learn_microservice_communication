### **Day 4: RPC & gRPC Fundamentals**

Today, we move away from standard REST APIs (using JSON over HTTP/1.1) and look at the industry standard for fast, internal microservice communication: **gRPC**.

#### **1. What is RPC? (Remote Procedure Call)**

In a REST API, you think in terms of "Resources" (e.g., `GET /inventory/123`).
In RPC, you think in terms of "Actions" or "Functions." The goal of RPC is to make a network call look and feel exactly like calling a local function in your code, like `checkStock(item_id)`.

#### **2. Why gRPC over REST?**

gRPC was developed by Google and has become the standard for service-to-service communication. It solves a lot of the performance bottlenecks of REST:

- **Protocol Buffers (Protobuf) instead of JSON:** JSON is text-based, heavy, and requires your CPU to parse it. Protobuf serializes your data into a highly compressed **binary** format. It's much smaller over the network and incredibly fast to serialize/deserialize.
- **HTTP/2 instead of HTTP/1.1:** REST generally uses HTTP/1.1, which handles one request at a time per TCP connection. gRPC uses HTTP/2, which allows multiplexing (sending hundreds of requests simultaneously over a single connection) and supports bidirectional streaming.
- **Strict Contracts:** In REST, you just hope the other service sends the JSON structure you expect. In gRPC, you define a strict contract using a `.proto` file. Both the Go client and the Go server generate their code from this exact same file. If the contract changes, the code won't compile, saving you from nasty runtime bugs.

#### **3. The Protobuf Contract (`.proto`)**

This is the heart of gRPC. It is language-agnostic. You write this file once, and you can generate Go code for your Order service and Python code for your Inventory service, and they will understand each other perfectly.

Here is what a simple contract for our Inventory system looks like:

```protobuf
// inventory.proto
syntax = "proto3";

// This tells the compiler where to put the generated Go code
option go_package = "./pb";

// The Request message (like the JSON body)
message StockRequest {
  string item_id = 1;
}

// The Response message
message StockResponse {
  bool in_stock = 1;
}

// The actual Service definition (the functions you can call)
service InventoryService {
  rpc CheckStock (StockRequest) returns (StockResponse);
}
```

_(Note: The `= 1` is not the value; it's a unique tag used for the binary compression)._

---

### **Actionable Task for Today**

Since you are using Go, we need to get your machine ready to compile `.proto` files into Go code for tomorrow.

1. **Install the Protobuf Compiler (`protoc`):**
   - Mac (Homebrew): `brew install protobuf`
   - Windows: Download the zip from the [official releases](https://github.com/protocolbuffers/protobuf/releases), extract it, and add it to your PATH.
   - Linux: `sudo apt install -y protobuf-compiler`
2. **Install the Go plugins for `protoc`:**
   Run these two commands in your terminal:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```
   _(Make sure your Go `bin` directory is in your system PATH so your terminal can find these plugins!)_
3. **Write the contract:** In your `day3-sync` folder (or make a new `day4-grpc` folder), create a file named `inventory.proto` and paste the protobuf code from the section above into it.

---

### **Day 4 Revision Question**

If gRPC is so much faster, uses less bandwidth, and has strict type-safety, why don't we use gRPC for _everything_? For example, why does your browser still use standard HTTP/REST to talk to a website's backend, instead of gRPC?

**Think about how browsers work versus how backend servers work. Let me know your answer, and we'll move on to Day 5, where we will actually write the gRPC code in Go!**

**Answer:** it's about how JavaScript works inside the browser. gRPC relies heavily on a feature of HTTP/2 called **"Trailers"** (which are basically HTTP headers sent at the very _end_ of a response, instead of the beginning). Browsers simply do not expose a way for JavaScript (like `fetch()` or `XMLHttpRequest`) to read these trailing headers. Because JS can't read the trailers, it can't understand if the gRPC call succeeded or failed.

_(Side note: If you ever absolutely need your web frontend to talk to a gRPC backend, you have to use a tool called `gRPC-Web`, which acts as a translator proxy between the browser and the server)._
