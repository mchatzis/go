package grpc

import (
	"context"
	"time"

	prod_grpc "github.com/mchatzis/go/producer/pkg/grpc"
	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/mchatzis/go/producer/pkg/sqlc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logger = logging.GetLogger()

func SendTasks(taskChan <-chan *sqlc.Task) {
	conn, err := grpc.NewClient("consumer:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Failed to connect to GRPC server: %v", err)
	}
	defer conn.Close()

	client := prod_grpc.NewTaskServiceClient(conn)

	logger.Info("Grpc client is running on port 50051...")
	logger.Info("Successfully connected to GRPC server...")

	for i := 0; i < 3; i++ {
		go func() {
			for task := range taskChan {
				sendTask(client, task)
			}
		}()
	}
	select {}
}

func sendTask(client prod_grpc.TaskServiceClient, task *sqlc.Task) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
		logger.Errorf("Error sending task: %v with error %v", task.ID, err)
		return
	}
	logger.Debugf("Grpc sent task: %v", task.ID)
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
