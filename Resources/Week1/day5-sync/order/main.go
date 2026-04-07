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
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
