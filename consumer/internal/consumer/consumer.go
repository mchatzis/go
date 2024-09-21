package consumer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/mchatzis/go/consumer/internal/grpc"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

const bufferSize = 1

func Consume(queries *sqlc.Queries) {
	var processedUnmatchedTasks sync.Map

	taskUpdateToProcessingStateChanIn := make(chan *sqlc.Task, bufferSize)
	taskUpdateToProcessingStateChanOut := make(chan *sqlc.Task, bufferSize)
	taskProcessChanIn := make(chan *sqlc.Task, bufferSize)
	taskProcessChanOut := make(chan *sqlc.Task, bufferSize)
	taskUpdateToDoneStateChanIn := make(chan *sqlc.Task, bufferSize)
	taskUpdateToDoneStateChanOut := make(chan *sqlc.Task, bufferSize)

	go grpc.ListenForTasks(taskProcessChanIn, taskUpdateToProcessingStateChanIn)
	for i := 0; i < 15; i++ {
		go processTask(taskProcessChanIn, taskProcessChanOut)
		go updateTaskState(taskUpdateToProcessingStateChanIn, taskUpdateToProcessingStateChanOut, sqlc.TaskStateProcessing, queries)
		go matchTasks(taskProcessChanOut, taskUpdateToProcessingStateChanOut, taskUpdateToDoneStateChanIn, &processedUnmatchedTasks)
		go updateTaskState(taskUpdateToDoneStateChanIn, taskUpdateToDoneStateChanOut, sqlc.TaskStateDone, queries)
		go logDoneTasks(taskUpdateToDoneStateChanOut)
	}
}

func matchTasks(taskChanIn chan *sqlc.Task, taskChanIn2 chan *sqlc.Task, taskChanOut chan *sqlc.Task, unmatchedTasks *sync.Map) {
	// Uses a map to match tasks incoming from the 'processing' and 'update-to-processing' channels.
	// Ensures each task has both been processed and updated in db to 'processing',
	// before forwarding to another db update to 'done' state.
	for {
		select {
		case task := <-taskChanIn:
			tryMatchTask(task, taskChanOut, unmatchedTasks)
		case task := <-taskChanIn2:
			tryMatchTask(task, taskChanOut, unmatchedTasks)
		}
	}
}

func tryMatchTask(task *sqlc.Task, taskChanOut chan *sqlc.Task, unmatchedTasks *sync.Map) {
	if _, exists := unmatchedTasks.Load(task.ID); exists {
		taskChanOut <- task
		unmatchedTasks.Delete(task.ID)
	} else {
		unmatchedTasks.Store(task.ID, *task)
	}
}

func processTask(taskChanIn chan *sqlc.Task, taskChanOut chan *sqlc.Task) {
	for task := range taskChanIn {
		time.Sleep(time.Duration(task.Value) * time.Millisecond)
		log.Printf("Processing task: %+v", task.ID)
		taskChanOut <- task
	}
}

func updateTaskState(taskChanIn chan *sqlc.Task, taskChanOut chan *sqlc.Task, state sqlc.TaskState, queries *sqlc.Queries) {
	for task := range taskChanIn {
		err := queries.UpdateTaskState(context.Background(), sqlc.UpdateTaskStateParams{
			State: state,
			ID:    task.ID,
		})
		if err != nil {
			log.Printf("Failed to update task state: %v", err)
		}
		log.Printf("Updated task to %v: %v", state, task.ID)
		taskChanOut <- task
	}
}

func logDoneTasks(taskChanIn chan *sqlc.Task) {
	for task := range taskChanIn {
		log.Printf("Done processing and updating task: %+v", task.ID)
	}
}
