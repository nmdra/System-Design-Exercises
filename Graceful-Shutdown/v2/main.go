package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(),
	}

	go func() {
		log.Println("[INFO] Server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Server failed: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	log.Println("[INFO] Interrupt signal received, initiating shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR] Graceful shutdown failed: %v", err)
	} else {
		log.Println("[INFO] Server shutdown complete")
	}

	wg.Wait()
	log.Println("[INFO] All background tasks completed, server exited cleanly")
}
