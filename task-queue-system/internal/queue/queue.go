package queue

import (
	"context"
	"log"

	"task-queue-system/config"

	"github.com/redis/go-redis/v9"
)

// Client is a package-level variable that holds the Redis client.
var Client *redis.Client

// TaskQueueName is the name of the list in Redis used as the task queue.
const TaskQueueName = "task_queue"

// ConnectQueue initializes the connection to the Redis server.
func ConnectQueue(cfg config.Config) {
	log.Println("Connecting to Redis...")
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	// Ping the Redis server to ensure a connection is established.
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection successful.")
	Client = client
}
