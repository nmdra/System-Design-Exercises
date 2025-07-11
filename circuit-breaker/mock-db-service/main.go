package main

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Name  string `json:"name,omitempty"`
	Error string `json:"error,omitempty"`
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		n := rand.Intn(10)

		switch {
		case n < 2:
			// Simulate internal server error
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Error: "simulated DB error"})
			slog.Error("DB returned error", "status", http.StatusInternalServerError)

		case n < 5:
			// Simulate slow response
			delay := time.Duration(2+rand.Intn(3)) * time.Second
			time.Sleep(delay)
			json.NewEncoder(w).Encode(Response{Name: "Alice (slow)"})
			slog.Warn("DB responded slowly", "delay", delay.String())

		default:
			// Normal response
			json.NewEncoder(w).Encode(Response{Name: "Alice"})
			slog.Info("DB responded successfully", "status", http.StatusOK)
		}
	})

	slog.Info("Mock DB service running", "url", "http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
