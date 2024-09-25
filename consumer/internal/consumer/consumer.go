package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/mchatzis/go/consumer/internal/grpc"
	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

const bufferSize = 50
const rateLimiterMultiplier = 500

var logger = logging.GetLogger()

/*
Updating the state to 'Processing' in DB happens concurrently with actually processing
the tasks. This is to avoid the I/O performance bottleneck of a sequential design, which
depends on DB access latency.
*/
func HandleIncomingTasks(queries *sqlc.Queries) {
	taskIncomingChan := make(chan *sqlc.Task)
	go grpc.ListenForTasks(taskIncomingChan)

	taskUpdateInDbToProcessingChanIn := make(chan *sqlc.Task, bufferSize)
	taskUpdateInDbToProcessingChanOut := make(chan *sqlc.Task, bufferSize)
	taskProcessChanIn := make(chan *sqlc.Task, bufferSize)
	taskProcessChanOut := make(chan *sqlc.Task, bufferSize)
	taskUpdateInDbToDoneChanIn := make(chan *sqlc.Task, bufferSize)
	taskUpdateInDbToDoneChanOut := make(chan *sqlc.Task, bufferSize)

	var unmatchedTasks sync.Map
	for i := 0; i < 15; i++ {
		go distributeIncomingTasks(taskIncomingChan, taskProcessChanIn, taskUpdateInDbToProcessingChanIn)
		go processTasks(taskProcessChanIn, taskProcessChanOut)
		go updateTasksStateInDb(taskUpdateInDbToProcessingChanIn, taskUpdateInDbToProcessingChanOut, sqlc.TaskStateProcessing, queries)
		go recombineChannels(taskProcessChanOut, taskUpdateInDbToProcessingChanOut, taskUpdateInDbToDoneChanIn, &unmatchedTasks)
		go updateTasksStateInDb(taskUpdateInDbToDoneChanIn, taskUpdateInDbToDoneChanOut, sqlc.TaskStateDone, queries)
		go logDoneTasks(taskUpdateInDbToDoneChanOut)
	}
}

func distributeIncomingTasks(taskChanIn <-chan *sqlc.Task, taskChanOut chan<- *sqlc.Task, taskChanOut2 chan<- *sqlc.Task) {
	for task := range taskChanIn {
		task.State = sqlc.TaskStateProcessing
		taskChanOut <- task
		taskChanOut2 <- task
		time.Sleep(time.Duration(rateLimiterMultiplier) * time.Millisecond)
	}
}

func processTasks(taskChanIn <-chan *sqlc.Task, taskChanOut chan<- *sqlc.Task) {
	for task := range taskChanIn {
		logger.Debugf("Processing task: %+v", task.ID)
		time.Sleep(time.Duration(task.Value) * time.Millisecond)
		task.State = sqlc.TaskStateDone
		taskChanOut <- task
	}
}

func updateTasksStateInDb(taskChanIn <-chan *sqlc.Task, taskChanOut chan<- *sqlc.Task, state sqlc.TaskState, queries *sqlc.Queries) {
	for task := range taskChanIn {
		err := queries.UpdateTaskState(context.Background(), sqlc.UpdateTaskStateParams{
			State: state,
			ID:    task.ID,
		})
		if err != nil {
			logger.Errorf("Failed to update task state: %v", err)
		}
		logger.Debugf("Updated task to %v: %v", state, task.ID)
		taskChanOut <- task
	}
}

func recombineChannels(taskChanIn <-chan *sqlc.Task, taskChanIn2 <-chan *sqlc.Task, taskChanOut chan<- *sqlc.Task, unmatchedTasks *sync.Map) {
	// Uses a map to match tasks incoming from the 'processing' and 'update-to-processing' channels.
	// Matched tasks are then sent to get updated to 'done'.
	// Ensures each task has both been processed and updated in db to 'processing',
	// 	before forwarding to another db update to 'done' state.
	for {
		select {
		case task := <-taskChanIn:
			tryMatchTask(task, taskChanOut, unmatchedTasks)
		case task := <-taskChanIn2:
			tryMatchTask(task, taskChanOut, unmatchedTasks)
		}
	}
}

func tryMatchTask(task *sqlc.Task, taskChanOut chan<- *sqlc.Task, unmatchedTasks *sync.Map) {
	if _, exists := unmatchedTasks.Load(task.ID); exists {
		taskChanOut <- task
		unmatchedTasks.Delete(task.ID)
	} else {
		unmatchedTasks.Store(task.ID, task)
	}
}

func logDoneTasks(taskChanIn <-chan *sqlc.Task) {
	for task := range taskChanIn {
		logger.Debugf("Done processing and updating task: %+v", task.ID)
	}
}
