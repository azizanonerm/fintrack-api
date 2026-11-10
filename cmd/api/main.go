package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/azizanonerm/fintrack-api/internal/config"
	"github.com/azizanonerm/fintrack-api/internal/repository"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := repository.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("database connection established")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	address := ":8080"
	log.Printf("server running at http://localhost%s", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal(err)
	}
}
