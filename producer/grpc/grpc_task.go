package grpc

import (
	"context"
	"log"
	"time"

	"github.com/mchatzis/go/producer/db/sqlc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TaskError struct {
	Task sqlc.Task
	Err  error
}

func SendTask(client TaskServiceClient, task sqlc.Task, errorChan chan<- TaskError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	grpcState := mapTaskState(task.State)
	_, err := client.SendTask(ctx, &Task{
		Id:             task.ID,
		Type:           task.Type,
		Value:          task.Value,
		State:          grpcState,
		CreationTime:   task.Creationtime,
		LastUpdateTime: task.Lastupdatetime,
	})
	if err != nil {
		log.Printf("Error sending task: %v", err)
		errorChan <- struct {
			Task sqlc.Task
			Err  error
		}{Task: task, Err: err}
		return
	}
}

func SendTasks(taskChan <-chan sqlc.Task, errorChan chan<- TaskError) {
	conn, err := grpc.NewClient("consumer:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := NewTaskServiceClient(conn)

	for task := range taskChan {
		SendTask(client, task, errorChan)
	}
}

func mapTaskState(state sqlc.TaskState) TaskState {
	switch state {
	case sqlc.TaskStatePending:
		return TaskState_PENDING
	case sqlc.TaskStateInProgress:
		return TaskState_IN_PROGRESS
	case sqlc.TaskStateCompleted:
		return TaskState_COMPLETED
	case sqlc.TaskStateFailed:
		return TaskState_FAILED
	default:
		panic("Task state failed to convert to grpc.TaskState")
	}
}
