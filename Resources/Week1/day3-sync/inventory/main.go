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
