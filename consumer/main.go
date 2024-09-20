package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	db "github.com/mchatzis/go/producer/db"
	"github.com/mchatzis/go/producer/db/sqlc"
)

const BufferSize = 100

func main() {
	connection := &db.Connection{}
	if err := connection.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer connection.Close()
	queries := sqlc.New(connection.Pool)

	var doneUnmatchedTasks sync.Map
	taskUpdateToInProgressChanIn := make(chan *sqlc.Task, BufferSize)
	taskUpdateToInProgressChanOut := make(chan *sqlc.Task, BufferSize)
	taskProcessChanIn := make(chan *sqlc.Task, BufferSize)
	taskProcessChanOut := make(chan *sqlc.Task, BufferSize)
	taskUpdateToCompletedChanIn := make(chan *sqlc.Task, BufferSize)
	taskUpdateToCompletedChanOut := make(chan *sqlc.Task, BufferSize)

	go ListenForTasks(taskProcessChanIn, taskUpdateToInProgressChanIn)
	for i := 0; i < 15; i++ {
		go processTask(taskProcessChanIn, taskProcessChanOut)
		go UpdateTaskState(taskUpdateToInProgressChanIn, taskUpdateToInProgressChanOut, sqlc.TaskStateInProgress, queries)
		go regroupTasks(taskProcessChanOut, taskUpdateToInProgressChanOut, taskUpdateToCompletedChanIn, &doneUnmatchedTasks, queries)
		go UpdateTaskState(taskUpdateToCompletedChanIn, taskUpdateToCompletedChanOut, sqlc.TaskStateCompleted, queries)
	}

	select {}
}

func regroupTasks(taskChanIn chan *sqlc.Task, taskChanIn2 chan *sqlc.Task, taskChanOut chan *sqlc.Task, doneUnmatchedTasks *sync.Map, queries *sqlc.Queries) {
	// Uses a map to match tasks incoming from the processing and update-to-in-progress channels.
	// Ensures task has both been processed and updated in db, before forwarding to another db update.
	for {
		select {
		case task := <-taskChanIn:
			forwardIfInMap(task, taskChanOut, doneUnmatchedTasks)
		case task := <-taskChanIn2:
			forwardIfInMap(task, taskChanOut, doneUnmatchedTasks)
		}
	}
}

func forwardIfInMap(task *sqlc.Task, taskChanOut chan *sqlc.Task, doneUnmatchedTasks *sync.Map) {
	if _, exists := doneUnmatchedTasks.Load(task.ID); exists {
		taskChanOut <- task
		doneUnmatchedTasks.Delete(task.ID)
	} else {
		doneUnmatchedTasks.Store(task.ID, *task)
	}
}

func processTask(taskChanIn chan *sqlc.Task, taskChanOut chan *sqlc.Task) {
	for task := range taskChanIn {
		log.Printf("Processing task: %+v", *task)
		time.Sleep(time.Duration(task.Value) * time.Millisecond)
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
	}
}
