package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/mchatzis/go/producer/pkg/base"
	prod_grpc "github.com/mchatzis/go/producer/pkg/grpc"
	"github.com/mchatzis/go/producer/pkg/logging"
)

const rateLimiterMultiplier int = 10

var logger = logging.GetLogger()

type server struct {
	prod_grpc.UnimplementedTaskServiceServer
	taskChan chan *base.Task
}

func (s *server) SendTask(ctx context.Context, grpc_task *prod_grpc.Task) (*emptypb.Empty, error) {
	task := &base.Task{}
	task.FromGRPCTask(grpc_task)

	s.taskChan <- task

	return &emptypb.Empty{}, nil
}

func ListenForTasks(taskChanOut chan *base.Task) {
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
