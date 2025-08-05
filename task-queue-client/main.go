package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "task-queue-client/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// serverAddresss is the address of the gRPC server.
	serverAddress = "localhost:50051"
)

func main() {
	// --- Command-Line Flag Definitions ---
	// Define a flag for submitting a task.
	// Example: go run main.go -submit -type="send_email" -payload="{\"to\":\"user@example.com\"}"
	submit := flag.Bool("submit", false, "Submit a new task")
	taskType := flag.String("type", "", "The type of the task to submit")
	taskPayload := flag.String("payload", "{}", "The JSON payload for the task")

	// Define a flag for checking the status of a task.
	// Example: go run main.go -status -id="your-task-id-here"
	status := flag.Bool("status", false, "Check the status of a task")
	taskID := flag.String("id", "", "The ID of the task to check")

	// Parse the command-line flags provided by the user.
	flag.Parse()

	// --- gRPC Connection Setup ---
	// Set up a connection to the gRPC server.
	// We use insecure credentials because we are running in a trusted local environment.
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	// Defer closing the connection until the main function returns.
	defer conn.Close()

	// Create a new TaskService client from the connection.
	client := pb.NewTaskServiceClient(conn)

	// --- Logic to Call gRPC Methods ---
	// Check which flag was used and call the appropriate function.
	if *submit {
		// If -submit flag is used, call the submitNewTask function.
		submitNewTask(client, *taskType, *taskPayload)
	} else if *status {
		// If -status flag is used, call the checkTaskStatus function.
		checkTaskStatus(client, *taskID)
	} else {
		// If no valid flag is provided, print usage instructions.
		log.Println("No action specified. Use -submit or -status flag.")
		flag.Usage()
	}
}

// submitNewTask calls the SubmitTask RPC.
func submitNewTask(client pb.TaskServiceClient, taskType, taskPayload string) {
	if taskType == "" {
		log.Fatal("Task type is required when submitting a task. Use the -type flag.")
	}

	log.Printf("Submitting task of type '%s' ...", taskType)

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Create the gRPC request message from the flags.
	req := &pb.SubmitTaskRequest{
		Type:    taskType,
		Payload: taskPayload,
	}

	// Call the remote SubmitTask method.
	res, err := client.SubmitTask(ctx, req)
	if err != nil {
		log.Fatalf("Could not submit task: %v", err)
	}

	// Print the server's response.
	log.Printf("✅ Task submitted successfully")
	log.Printf("   Task ID: %s", res.GetTaskId())
}

// checkTaskStatus calls the GetTaskStatus RPC.
func checkTaskStatus(client pb.TaskServiceClient, taskID string) {
	if taskID == "" {
		log.Fatal("Task ID is required when checking status. Use the -id flag.")
	}

	log.Printf("Checking status for task ID '%s' ...", taskID)

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Create the gRPC request message.
	req := &pb.GetTaskStatusRequest{TaskId: taskID}

	// Call the remote GetTaskStatus method.
	res, err := client.GetTaskStatus(ctx, req)
	if err != nil {
		log.Fatalf("Could not get task status: %v", err)
	}

	// Print the server's response.
	log.Printf("✅ Status received successfully")
	log.Printf("   Task ID: %s", res.GetTaskId())
	log.Printf("   Status: %s", res.GetStatus())
	log.Printf("   Result: %s", res.GetResult())
}
