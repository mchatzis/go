package producer

import (
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchatzis/go/producer/internal/db"
	"github.com/mchatzis/go/producer/internal/grpc"
	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

const MaxBacklog int = 50

var logger = logging.GetLogger()

func Produce(pool *pgxpool.Pool) {
	queries := sqlc.New(pool)
	time.Sleep(time.Second) //Wait a sec for postgres to spin up, should apply a healthcheck instead in the future
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	saveTaskChan := make(chan sqlc.Task, MaxBacklog)
	sendTaskChan := make(chan sqlc.Task, MaxBacklog)

	for i := 0; i < 3; i++ {
		go saveTasks(queries, saveTaskChan)
	}
	go grpc.SendTasks(sendTaskChan)

	logger.Info("Generating and sending tasks...")
	for i := 1; ; i++ {
		task := generateTask(r, i)
		sendTaskChan <- task
		saveTaskChan <- task
	}
}

func generateTask(r *rand.Rand, id int) sqlc.Task {
	if id <= 0 {
		panic("ID field must be positive")
	}

	taskType := r.Intn(10)
	taskValue := r.Intn(100)

	task := sqlc.Task{
		ID:             int32(id),
		Type:           int32(taskType),
		Value:          int32(taskValue),
		State:          sqlc.TaskStatePending,
		Creationtime:   float64(time.Now().Unix()),
		Lastupdatetime: float64(time.Now().Unix()),
	}
	logger.Debugf("Generated task: %v", task.ID)
	return task
}

func saveTasks(queries *sqlc.Queries, taskChan <-chan sqlc.Task) {
	for task := range taskChan {
		err := db.CreateTask(queries, task)
		if err != nil {
			logger.Errorf("Failed to save task %v with error: %v", task.ID, err)
		}
		logger.Debugf("Created in db pending task: %v", task.ID)
	}
}
