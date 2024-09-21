package grpc

import (
	"context"
	"log"
	"time"

	prod_grpc "github.com/mchatzis/go/producer/pkg/grpc"
	"github.com/mchatzis/go/producer/pkg/sqlc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendTasks(taskChan <-chan sqlc.Task) {
	conn, err := grpc.NewClient("consumer:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := prod_grpc.NewTaskServiceClient(conn)

	for task := range taskChan {
		sendTask(client, task)
	}
}

func sendTask(client prod_grpc.TaskServiceClient, task sqlc.Task) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	grpcState := mapTaskState(task.State)
	_, err := client.SendTask(ctx, &prod_grpc.Task{
		Id:             task.ID,
		Type:           task.Type,
		Value:          task.Value,
		State:          grpcState,
		CreationTime:   task.Creationtime,
		LastUpdateTime: task.Lastupdatetime,
	})
	if err != nil {
		log.Printf("Error sending task: %v with error %v", task.ID, err)
		return
	}
}

func mapTaskState(state sqlc.TaskState) prod_grpc.TaskState {
	switch state {
	case sqlc.TaskStatePending:
		return prod_grpc.TaskState_PENDING
	case sqlc.TaskStateProcessing:
		return prod_grpc.TaskState_PROCESSING
	case sqlc.TaskStateDone:
		return prod_grpc.TaskState_DONE
	case sqlc.TaskStateFailed:
		return prod_grpc.TaskState_FAILED
	default:
		panic("Task state failed to convert to grpc.TaskState")
	}
}
