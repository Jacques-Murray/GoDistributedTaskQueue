package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"task-queue-system/internal/database"
	"task-queue-system/internal/queue"

	"github.com/google/uuid"
)

// Start begins the worker process of listening for and processing tasks.
func Start() {
	log.Println("Worker started. Waiting for tasks...")

	for {
		// Block and wait for a task to be pushed to the queue.
		// BRPOP is a blocking right-pop operation.
		result, err := queue.Client.BRPop(context.Background(), 0, queue.TaskQueueName).Result()
		if err != nil {
			log.Printf("Error dequeueing task: %v", err)
			continue // Continue to the next iteration on error.
		}

		// result is a slice where result[0] is the queue name and result[1] is the value.
		taskPayload := result[1]

		var payload map[string]string
		if err := json.Unmarshal([]byte(taskPayload), &payload); err != nil {
			log.Printf("Error unmarshalling task payload: %v", err)
			continue
		}

		taskIDStr := payload["task_id"]
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			log.Printf("Invalid TaskID received from queue: %v", err)
			continue
		}

		// Process the task in a separate goroutine to allow concurrent processing.
		go processTask(taskID)
	}
}

// processTask fetches a task by ID, executes it, and updates its status.
func processTask(taskID uuid.UUID) {
	log.Printf("Processing task: %s", taskID)

	var task database.Task
	// Fetch the full task details from the database.
	if err := database.DB.First(&task, taskID).Error; err != nil {
		log.Printf("Error fetching task %s from DB: %v", taskID, err)
		return
	}

	// Update task status to RUNNING.
	database.DB.Model(&task).Update("status", "RUNNING")

	// --- Simulate Task Execution ---
	// In a real application, you would have a switch statement on task.Type
	// to call the appropriate function.
	log.Printf("Executing task type '%s' with payload: %s", task.Type, task.Payload)
	time.Sleep(5 * time.Second) // Simulate work being done.
	taskResult := `{"status": "success", "message": "Task completed successfully"}`
	// --- End Simulation ---

	// Update task status to COMPLETED and store the result.
	database.DB.Model(&task).Updates(database.Task{Status: "COMPLETED", Result: taskResult})

	log.Printf("Finished processing task: %s", taskID)
}
