package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	address := ":8080"
	log.Printf("server running at http://localhost%s", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal(err)
	}
}
