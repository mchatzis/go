package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	prod_grpc "github.com/mchatzis/go/producer/pkg/grpc"
	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

const rateLimiterMultiplier int = 10

var logger = logging.GetLogger()

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

func ListenForTasks(taskChanOut chan *sqlc.Task) {
	logger.Info("Opening tcp connection...")
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatalf("Failed to listen tcp on port 50051 with error: %v", err)
	}

	server := &server{
		taskChan: taskChanOut,
	}
	grpcServer := grpc.NewServer()
	prod_grpc.RegisterTaskServiceServer(grpcServer, server)

	logger.Info("Listening for tasks on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatalf("Grpc server failed: %v", err)
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
		logger.Fatalf("Unexpected TaskState value: %v", state)
		return ""
	}
}
