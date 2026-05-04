package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"push-n-pray/internal/backend"
)

func main() {
	addr := os.Getenv("BACKEND_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("backend listening on %s", addr)
	if err := backend.Run(ctx, addr); err != nil {
		log.Fatalf("backend stopped with error: %v", err)
	}
}
