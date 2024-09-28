package grpc

import (
	"context"
	"time"

	"github.com/mchatzis/go/producer/pkg/base"
	prod_grpc "github.com/mchatzis/go/producer/pkg/grpc"
	"github.com/mchatzis/go/producer/pkg/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logger = logging.GetLogger()

func SendTasks(taskChan <-chan *base.Task) {
	conn, err := grpc.NewClient("consumer:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Failed to connect to GRPC server: %v", err)
	}
	defer conn.Close()

	client := prod_grpc.NewTaskServiceClient(conn)

	logger.Info("Grpc client is running on port 50051...")
	logger.Info("Successfully connected to GRPC server...")

	for task := range taskChan {
		sendTask(client, task)
	}
}

func sendTask(client prod_grpc.TaskServiceClient, task *base.Task) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := client.SendTask(ctx, task.ToGRPCTask())
	if err != nil {
		logger.Errorf("Error sending task: %v with error %v", task.ID, err)
		return
	}
	logger.Debugf("Grpc sent task: %v", task.ID)
}
