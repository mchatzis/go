package consumer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/mchatzis/go/consumer/internal/grpc"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

const BufferSize = 1

func consume(queries *sqlc.Queries) {
	var processedUnmatchedTasks sync.Map

	taskUpdateToProcessingStateChanIn := make(chan *sqlc.Task, BufferSize)
	taskUpdateToProcessingStateChanOut := make(chan *sqlc.Task, BufferSize)
	taskProcessChanIn := make(chan *sqlc.Task, BufferSize)
	taskProcessChanOut := make(chan *sqlc.Task, BufferSize)
	taskUpdateToDoneStateChanIn := make(chan *sqlc.Task, BufferSize)
	taskUpdateToDoneStateChanOut := make(chan *sqlc.Task, BufferSize)

	go grpc.ListenForTasks(taskProcessChanIn, taskUpdateToProcessingStateChanIn)
	for i := 0; i < 15; i++ {
		go processTask(taskProcessChanIn, taskProcessChanOut)
		go UpdateTaskState(taskUpdateToProcessingStateChanIn, taskUpdateToProcessingStateChanOut, sqlc.TaskStateProcessing, queries)
		go matchTasks(taskProcessChanOut, taskUpdateToProcessingStateChanOut, taskUpdateToDoneStateChanIn, &processedUnmatchedTasks)
		go UpdateTaskState(taskUpdateToDoneStateChanIn, taskUpdateToDoneStateChanOut, sqlc.TaskStateDone, queries)
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
		taskChanOut <- task
	}
}

func logDoneTasks(taskChanIn chan *sqlc.Task) {
	for task := range taskChanIn {
		log.Printf("Done processing and updating task: %+v", task.ID)
	}
}
