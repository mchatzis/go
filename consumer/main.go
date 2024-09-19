package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	db "github.com/mchatzis/go/producer/db"
	"github.com/mchatzis/go/producer/db/sqlc"
)

func main() {
	connection := &db.Connection{}
	if err := connection.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer connection.Close()
	queries := sqlc.New(connection.Pool)

	taskUpdateToInProgressChan := make(chan *sqlc.Task, 100)
	go UpdateTaskStateToInProgress(taskUpdateToInProgressChan, queries)

	taskProcessChan := make(chan *sqlc.Task, 100)
	for i := 0; i < 15; i++ {
		go processTask(taskProcessChan)
	}

	taskListenChan := make(chan *sqlc.Task, 100)
	go ListenForTasks(taskListenChan)

	for task := range taskListenChan {
		taskProcessChan <- task
		taskUpdateToInProgressChan <- task
	}

	select {}
}

func processTask(taskChan chan *sqlc.Task) {
	for task := range taskChan {
		log.Printf("Processing task: %+v", *task)
		time.Sleep(time.Duration(task.Value) * time.Millisecond)
		log.Printf("Finished processing task: %+v", *task)
	}
}

func UpdateTaskStateToInProgress(taskChan chan *sqlc.Task, queries *sqlc.Queries) {
	for task := range taskChan {
		err := queries.UpdateTaskState(context.Background(), sqlc.UpdateTaskStateParams{
			State: sqlc.TaskStateInProgress,
			ID:    task.ID,
		})
		if err != nil {
			log.Printf("Failed to update task state: %v", err)
		}
		log.Printf("Updated task state to in_progress for task: %+v", *task)
	}
}
