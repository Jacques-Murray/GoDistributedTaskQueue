package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration values.
type Config struct {
	DatabaseURL string `envconfig:"POSTGRES_HOSTNAME"`
	RedisURL    string `envconfig:"REDIS_URL"`
	GRPCPort    string `envconfig:"GRPC_PORT"`
}

// LoadConfig load configuration from a .env file and environment variables.
func LoadConfig() Config {
	// Load values from .env file for local development.
	// In production, these should be set directly as environment variables.
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
