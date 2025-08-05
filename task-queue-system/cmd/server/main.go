package main

import (
	"log"
	"net"

	"task-queue-system/api"
	"task-queue-system/config"
	"task-queue-system/internal/database"
	"task-queue-system/internal/queue"
	pb "task-queue-system/proto"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting API Server...")

	// 1. Load configuration from environment variables.
	cfg := config.LoadConfig()
	log.Println("Configuration loaded.")

	// 2. Connect to the database.
	database.ConnectDatabase(cfg)

	// 3. Connect to the Redis queue.
	queue.ConnectQueue(cfg)

	// 4. Set up a TCP listener on the configured gRPC port.
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.GRPCPort, err)
	}
	log.Printf("Listening on port %s", cfg.GRPCPort)

	// 5. Create a new gRPC server instance.
	s := grpc.NewServer()

	// 6. Register our API server implementation with the gRPC server.
	pb.RegisterTaskServiceServer(s, &api.Server{})

	// 7. Start the gRPC server.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
