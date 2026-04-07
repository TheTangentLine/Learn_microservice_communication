package main

import (
	"encoding/json"
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
