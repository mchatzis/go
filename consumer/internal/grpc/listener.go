package grpc

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	prod_grpc "github.com/mchatzis/go/producer/pkg/grpc"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

const rateLimiterMultiplier int = 1

type server struct {
	prod_grpc.UnimplementedTaskServiceServer
	taskChan chan *sqlc.Task
}

func (s *server) SendTask(ctx context.Context, grpc_task *prod_grpc.Task) (*emptypb.Empty, error) {
	task := sqlc.Task{
		ID:             grpc_task.Id,
		Type:           grpc_task.Type,
		Value:          grpc_task.Value,
		State:          mapGrpcToSqlc(grpc_task.State),
		Creationtime:   grpc_task.CreationTime,
		Lastupdatetime: grpc_task.LastUpdateTime,
	}

	s.taskChan <- &task

	return &emptypb.Empty{}, nil
}

func ListenForTasks(taskChanOut chan *sqlc.Task, taskChanOut2 chan *sqlc.Task) {
	taskListenChan := make(chan *sqlc.Task)
	go listen(taskListenChan)

	for task := range taskListenChan {
		time.Sleep(time.Duration(rateLimiterMultiplier) * time.Millisecond)
		log.Printf("Received task: %+v", task.ID)
		taskChanOut <- task
		taskChanOut2 <- task
	}
}

func listen(taskChan chan *sqlc.Task) {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := &server{
		taskChan: taskChan,
	}
	grpcServer := grpc.NewServer()
	prod_grpc.RegisterTaskServiceServer(grpcServer, server)

	log.Println("Server is running on port 8082...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func mapGrpcToSqlc(state prod_grpc.TaskState) sqlc.TaskState {
	switch state {
	case prod_grpc.TaskState_PENDING:
		return sqlc.TaskStatePending
	case prod_grpc.TaskState_PROCESSING:
		return sqlc.TaskStateProcessing
	case prod_grpc.TaskState_DONE:
		return sqlc.TaskStateDone
	case prod_grpc.TaskState_FAILED:
		return sqlc.TaskStateFailed
	default:
		log.Fatalf("Unexpected TaskState value: %v", state)
		return ""
	}
}
