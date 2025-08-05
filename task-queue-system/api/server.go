package api

import (
	"context"
	"encoding/json"
	"log"

	"task-queue-system/internal/database"
	"task-queue-system/internal/queue"
	pb "task-queue-system/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is the struct that implements the gRPC TaskService.
type Server struct {
	pb.UnimplementedTaskServiceServer
}

// SubmitTask is the gRPC method for adding a new task to the queue.
func (s *Server) SubmitTask(ctx context.Context, req *pb.SubmitTaskRequest) (*pb.SubmitTaskResponse, error) {
	log.Printf("Received SubmitTask request: Type=%s", req.Type)

	// Before processing, validate that the incoming payload is a valid JSON string.
	// This prevents errors when saving to the jsonb database column.
	var js json.RawMessage
	if err := json.Unmarshal([]byte(req.Payload), &js); err != nil {
		log.Printf("Error: payload is not valid JSON: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "payload field must be a valid JSON string")
	}

	// Create a new task model from the request.
	task := database.Task{
		Type:     req.Type,
		Payload:  req.Payload,
		Priority: int(req.Priority),
		Status:   "PENDING",
	}

	// Save the new task to the database.
	result := database.DB.Create(&task)
	if result.Error != nil {
		log.Printf("Error creating task in DB: %v", result.Error)
		return nil, result.Error
	}

	// Prepare task data for the queue.
	queuePayload, err := json.Marshal(map[string]string{"task_id": task.ID.String()})
	if err != nil {
		log.Printf("Error marshalling queue payload: %v", err)
		return nil, err
	}

	// Enqueue the task ID for a worker to pick up.
	err = queue.Client.LPush(context.Background(), queue.TaskQueueName, string(queuePayload)).Err()
	if err != nil {
		log.Printf("Error enqueuing task: %v", err)
		return nil, err
	}

	log.Printf("Task %s enqueued successfully", task.ID)

	// Return the new task's ID in the response.
	return &pb.SubmitTaskResponse{TaskId: task.ID.String()}, nil
}

// GetTaskStatus is the gRPC method for checking a task's status.
func (s *Server) GetTaskStatus(ctx context.Context, req *pb.GetTaskStatusRequest) (*pb.GetTaskStatusResponse, error) {
	log.Printf("Received GetTaskStatus request for TaskID: %s", req.TaskId)

	var task database.Task
	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		log.Printf("Invalid TaskID format: %v", err)
		return nil, err
	}

	// Find the task in the database by its ID.
	result := database.DB.First(&task, taskID)
	if result.Error != nil {
		log.Printf("Error fetching task from DB: %v", result.Error)
		return nil, result.Error
	}

	// Return the task's status and result.
	return &pb.GetTaskStatusResponse{
		TaskId: task.ID.String(),
		Status: task.Status,
		Result: task.Result,
	}, nil
}
