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
	conn, err := grpc.NewClient("inventory:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
