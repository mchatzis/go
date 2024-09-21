package producer

import (
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchatzis/go/producer/internal/db"
	"github.com/mchatzis/go/producer/internal/grpc"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

const MaxBacklog int = 50

func Produce(pool *pgxpool.Pool) {
	queries := sqlc.New(pool)
	time.Sleep(time.Second)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	saveTaskChan := make(chan sqlc.Task, MaxBacklog)
	sendTaskChan := make(chan sqlc.Task, MaxBacklog)

	for i := 0; i < 3; i++ {
		go saveTasks(queries, saveTaskChan)
	}
	go grpc.SendTasks(sendTaskChan)

	for i := 1; ; i++ {
		task := generateTask(r, i)
		sendTaskChan <- task
		saveTaskChan <- task

		log.Print(task)
	}
}

func generateTask(r *rand.Rand, id int) sqlc.Task {
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

	return task
}

func saveTasks(queries *sqlc.Queries, taskChan <-chan sqlc.Task) {
	for task := range taskChan {
		err := db.CreateTask(queries, task)
		if err != nil {
			log.Printf("Failed to save task %v with error: %v", task.ID, err)
		}
	}
}
