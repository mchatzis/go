package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	producer_grpc "github.com/mchatzis/go/producer/grpc"
)

type TaskState string

const (
	TaskStatePending    TaskState = "pending"
	TaskStateInProgress TaskState = "in_progress"
	TaskStateCompleted  TaskState = "completed"
	TaskStateFailed     TaskState = "failed"
)

type Task struct {
	ID             int32
	Type           int32
	Value          int32
	State          TaskState
	Creationtime   float64
	Lastupdatetime float64
}

type server struct {
	producer_grpc.UnimplementedTaskServiceServer
	taskChan chan *Task
}

func (s *server) SendTask(ctx context.Context, task *producer_grpc.Task) (*emptypb.Empty, error) {
	consumerTask := Task{
		ID:             task.Id,
		Type:           task.Type,
		Value:          task.Value,
		State:          mapGrpcTaskState(task.State),
		Creationtime:   task.CreationTime,
		Lastupdatetime: task.LastUpdateTime,
	}

	log.Printf("Received task: %+v", consumerTask)

	s.taskChan <- &consumerTask

	return &emptypb.Empty{}, nil
}

func ListenForTasks(taskChan chan *Task) {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := &server{
		taskChan: taskChan,
	}
	grpcServer := grpc.NewServer()
	producer_grpc.RegisterTaskServiceServer(grpcServer, server)

	log.Println("Server is running on port 8082...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func mapGrpcTaskState(state producer_grpc.TaskState) TaskState {
	switch state {
	case producer_grpc.TaskState_PENDING:
		return TaskStatePending
	case producer_grpc.TaskState_IN_PROGRESS:
		return TaskStateInProgress
	case producer_grpc.TaskState_COMPLETED:
		return TaskStateCompleted
	case producer_grpc.TaskState_FAILED:
		return TaskStateFailed
	default:
		log.Fatalf("Unexpected TaskState value: %v", state)
		return ""
	}
}
