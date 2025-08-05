package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration values.
type Config struct {
	DatabaseURL string
	RedisURL    string
	GRPCPort    string
}

// LoadConfig load configuration from a .env file and environment variables.
func LoadConfig() Config {
	// Load values from .env file for local development. This is useful for development outside of Docker.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	// Construct the database URL from environment variables.
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := "db"
	dbName := os.Getenv("POSTGRES_DB")
	dbURL := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":5432/" + dbName + "?sslmode=disable"

	return Config{
		DatabaseURL: dbURL,
		RedisURL:    os.Getenv("REDIS_URL"),
		GRPCPort:    "50051",
	}
}
