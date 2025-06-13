package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type HelloResponse struct {
    Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    // Simulate concurrency task (if needed)
    go func() {
        log.Println("Handling request concurrently")
    }()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(HelloResponse{Message: "Hello, world!"})
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/hello", helloHandler)

    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  15 * time.Second,
    }

    log.Println("Server is running on http://localhost:8080")
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("could not start server: %v\n", err)
    }
}
