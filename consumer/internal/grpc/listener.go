package grpc

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/mchatzis/go/producer/pkg/base"
	prod_grpc "github.com/mchatzis/go/producer/pkg/grpc"
	"github.com/mchatzis/go/producer/pkg/logging"
)

var logger = logging.GetLogger()

type server struct {
	prod_grpc.UnimplementedTaskServiceServer
	taskChan chan *base.Task
}

func (s *server) SendTasks(stream prod_grpc.TaskService_SendTasksServer) error {
	for {
		received, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&emptypb.Empty{})
		}
		if err != nil {
			return err
		}

		task := &base.Task{}
		task.FromGRPCTask(received)

		s.taskChan <- task
	}

}

func ListenForTasks(taskChanOut chan *base.Task) {
	logger.Info("Opening tcp connection...")
	port := 50051
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
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
