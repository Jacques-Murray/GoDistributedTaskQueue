package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"task-queue-system/config"
	"task-queue-system/internal/database"
	"task-queue-system/internal/queue"
	"task-queue-system/internal/worker"
)

func main() {
	log.Println("Starting Worker...")

	// 1. Load configuration from environment variables.
	cfg := config.LoadConfig()
	log.Println("Configuration loaded.")

	// 2. Connect to the database.
	database.ConnectDatabase(cfg)

	// 3. Connect to the Redis queue.
	queue.ConnectQueue(cfg)

	// 4. Start the worker process in a separate goroutine.
	go worker.Start()

	// 5. Wait for a shutdown signal to gracefully exit.
	// This keeps the main function alive while the worker runs.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
}
