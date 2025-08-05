package database

import (
	"log"
	"task-queue-system/config"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a package-level variable that holds the database connection pool.
var DB *gorm.DB

// Task represents the GORM model for a task in the database.
// Each field is mapped to a column in the 'tasks' table.
type Task struct {
	// ID is the unique identifier for the task, a UUID generated on creation.
	ID uuid.UUID `gorm:"type:uuid,primary_key;"`

	// Type is a string that identifies the kind of task (e.g., "send_email").
	Type string

	// Payload stores the data required for the task to run, as a JSON object.
	Payload string `gorm:"type:jsonb"`

	// Priority determines the execution order. Higher numbers are higher priority.
	Priority int `gorm:"default:0"`

	// RetryCount tracks the number of times this task has been attempted.
	RetryCount int `gorm:"default:0"`

	// Status indicates the current state of the task (e.g., "PENDING", "RUNNING", "COMPLETED").
	Status string

	// Result stores the output of the task upon completion, as a JSON object.
	Result string `gorm:"type:jsonb"`

	// CreatedAt is the timestamp when the task was created.
	CreatedAt time.Time

	// UpdatedAt is the timestamp when the task was last updated.
	UpdatedAt time.Time
}

// BeforeCreate is a GORM hook that runs before a new Task record is created.
// It automatically generates a new UUID for the task's ID field.
func (task *Task) BeforeCreate(tx *gorm.DB) (err error) {
	task.ID = uuid.New()
	return
}

// ConnectDatabase initializes the connection to the PostgreSQL database.
// It uses the configuration provided to establish the connection and then
// runs AutoMigrate to ensure the 'tasks' table exists.
func ConnectDatabase(cfg config.Config) {
	log.Println("Connecting to database...")

	// Open a connection to the database using the DSN from the config.
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		// If the connection fails, log the error and exit the application.
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection successful.")

	// AutoMigrate will create or update the 'tasks' table to match the Task struct.
	log.Println("Running database migrations...")
	db.AutoMigrate(&Task{})
	log.Println("Database migrations complete.")

	// Assign the database connection pool to the package-level DB variable.
	DB = db
}
