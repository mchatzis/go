package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchatzis/go/producer/db"
	"github.com/mchatzis/go/producer/db/sqlc"
	"github.com/mchatzis/go/producer/grpc"
)

func Run(pool *pgxpool.Pool) {
	queries := sqlc.New(pool)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	taskChan := make(chan sqlc.Task)
	errorChan := make(chan struct {
		Task sqlc.Task
		Err  error
	}, 100)

	go saveTasks(queries, taskChan, errorChan)
	go grpc.SendTasks(taskChan, errorChan)

	var failedSaveTasks []sqlc.Task
	for i := 103; i < 108; i++ {
		task := generateTask(r, i)
		taskChan <- task

		select {
		case errInfo := <-errorChan:
			log.Printf("Error sending/saving task ID %d: %v\n", errInfo.Task.ID, errInfo.Err)
			failedSaveTasks = append(failedSaveTasks, errInfo.Task)
		default:
		}

		log.Print(task)
		time.Sleep(5 * time.Second)
	}
}

func saveTasks(queries *sqlc.Queries, taskChan <-chan sqlc.Task, errorChan chan<- struct {
	Task sqlc.Task
	Err  error
}) {
	for task := range taskChan {
		err := db.CreateTask(queries, task)
		if err != nil {
			errorChan <- struct {
				Task sqlc.Task
				Err  error
			}{Task: task, Err: err}
		}
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
