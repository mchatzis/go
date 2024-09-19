package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/mchatzis/go/producer/db/sqlc"
	producer_grpc "github.com/mchatzis/go/producer/grpc"
)

type server struct {
	producer_grpc.UnimplementedTaskServiceServer
	taskChan chan *sqlc.Task
}

func (s *server) SendTask(ctx context.Context, grpc_task *producer_grpc.Task) (*emptypb.Empty, error) {
	task := sqlc.Task{
		ID:             grpc_task.Id,
		Type:           grpc_task.Type,
		Value:          grpc_task.Value,
		State:          mapGrpcToSqlc(grpc_task.State),
		Creationtime:   grpc_task.CreationTime,
		Lastupdatetime: grpc_task.LastUpdateTime,
	}

	log.Printf("Received task: %+v", task)

	s.taskChan <- &task

	return &emptypb.Empty{}, nil
}

func ListenForTasks(taskChan chan *sqlc.Task) {
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

func mapGrpcToSqlc(state producer_grpc.TaskState) sqlc.TaskState {
	switch state {
	case producer_grpc.TaskState_PENDING:
		return sqlc.TaskStatePending
	case producer_grpc.TaskState_IN_PROGRESS:
		return sqlc.TaskStateInProgress
	case producer_grpc.TaskState_COMPLETED:
		return sqlc.TaskStateCompleted
	case producer_grpc.TaskState_FAILED:
		return sqlc.TaskStateFailed
	default:
		log.Fatalf("Unexpected TaskState value: %v", state)
		return ""
	}
}
