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

	taskUpdateToInProgressChanIn := make(chan *sqlc.Task, 100)
	taskUpdateToInProgressChanOut := make(chan *sqlc.Task, 100)
	taskProcessChanIn := make(chan *sqlc.Task, 100)
	taskProcessChanOut := make(chan *sqlc.Task, 100)
	for i := 0; i < 15; i++ {
		go processTask(taskProcessChanIn, taskProcessChanOut)
		go UpdateTaskState(taskUpdateToInProgressChanIn, taskUpdateToInProgressChanOut, sqlc.TaskStateInProgress, queries)
	}

	go ListenForTasks(taskProcessChanIn, taskUpdateToInProgressChanIn)

	select {}
}

func processTask(taskChanIn chan *sqlc.Task, taskChanOut chan *sqlc.Task) {
	for task := range taskChanIn {
		log.Printf("Processing task: %+v", *task)
		time.Sleep(time.Duration(task.Value) * time.Millisecond)
		log.Printf("Finished processing task: %+v", *task)
		taskChanOut <- task
	}
}

func UpdateTaskState(taskChanIn chan *sqlc.Task, taskChanOut chan *sqlc.Task, state sqlc.TaskState, queries *sqlc.Queries) {
	for task := range taskChanIn {
		err := queries.UpdateTaskState(context.Background(), sqlc.UpdateTaskStateParams{
			State: state,
			ID:    task.ID,
		})
		if err != nil {
			log.Printf("Failed to update task state: %v", err)
		}
		log.Printf("Updated task state to in_progress for task: %+v", *task)
		taskChanOut <- task
	}
}
