package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// The URL of our internal Order Service (using Docker's DNS)
	orderServiceURL, err := url.Parse("http://order:8080")
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
