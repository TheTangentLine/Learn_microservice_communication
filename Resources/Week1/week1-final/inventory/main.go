package main

import (
	"context"
	"log"
	"net"

	// Replace "day5-grpc" with your actual go module name if different
	"week1-final/pb"

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
