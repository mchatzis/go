package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/mchatzis/go/consumer/internal/grpc"
	"github.com/mchatzis/go/producer/pkg/base"
	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

const bufferSize = 50

var logger = logging.GetLogger()

/*
Updating the state to 'Processing' in DB happens concurrently with actually processing
the tasks. This is to avoid the I/O performance bottleneck of a sequential design, which
depends on DB access latency.
*/
func HandleIncomingTasks(queries *sqlc.Queries, rateLimitMultiplier int) {
	taskIncomingChan := make(chan *base.Task)
	go grpc.ListenForTasks(taskIncomingChan)

	taskUpdateInDbToProcessingChanIn := make(chan *base.Task, bufferSize)
	taskUpdateInDbToProcessingChanOut := make(chan *base.Task, bufferSize)
	taskProcessChanIn := make(chan *base.Task, bufferSize)
	taskProcessChanOut := make(chan *base.Task, bufferSize)
	taskUpdateInDbToDoneChanIn := make(chan *base.Task, bufferSize)
	taskUpdateInDbToDoneChanOut := make(chan *base.Task, bufferSize)

	var rateLimitDuration = time.Duration(rateLimitMultiplier) * time.Millisecond
	var unmatchedTasks sync.Map
	for i := 0; i < 15; i++ {
		go distributeIncomingTasksWithRateLimit(taskIncomingChan, taskProcessChanIn, taskUpdateInDbToProcessingChanIn, rateLimitDuration)
		go processTasks(taskProcessChanIn, taskProcessChanOut, pretendToProcess)
		go updateTasksStateInDb(taskUpdateInDbToProcessingChanIn, taskUpdateInDbToProcessingChanOut, base.TaskStateProcessing, queries)
		go recombineChannels(taskProcessChanOut, taskUpdateInDbToProcessingChanOut, taskUpdateInDbToDoneChanIn, &unmatchedTasks)
		go updateTasksStateInDb(taskUpdateInDbToDoneChanIn, taskUpdateInDbToDoneChanOut, base.TaskStateDone, queries)
		go logDoneTasks(taskUpdateInDbToDoneChanOut)
	}
}

func distributeIncomingTasksWithRateLimit(taskChanIn <-chan *base.Task, taskChanOut chan<- *base.Task, taskChanOut2 chan<- *base.Task, rateLimitDuration time.Duration) {
	for task := range taskChanIn {
		task.State = base.TaskStateProcessing
		taskChanOut <- task
		taskChanOut2 <- task
		time.Sleep(rateLimitDuration)
	}
}

func processTasks(taskChanIn <-chan *base.Task, taskChanOut chan<- *base.Task, process func(*base.Task) error) {
	for task := range taskChanIn {
		logger.Debugf("Processing task: %+v", task.ID)
		err := process(task)
		if err != nil {
			task.State = base.TaskStateFailed
			logger.Infof("Task %v failed processing with error: %v", task.ID, err)
		} else {
			task.State = base.TaskStateDone
			taskChanOut <- task
		}
	}
}

func pretendToProcess(task *base.Task) error {
	time.Sleep(time.Duration(task.Value) * time.Millisecond)
	return nil
}

func updateTasksStateInDb(taskChanIn <-chan *base.Task, taskChanOut chan<- *base.Task, state base.TaskState, queries *sqlc.Queries) {
	for task := range taskChanIn {
		sqlcTask := task.ToSQLCTask()
		err := queries.UpdateTaskState(context.Background(), sqlc.UpdateTaskStateParams{
			State: sqlcTask.State,
			ID:    sqlcTask.ID,
		})
		if err != nil {
			logger.Errorf("Failed to update task state: %v", err)
		}
		logger.Debugf("Updated task to %v: %v", state, task.ID)
		taskChanOut <- task
	}
}

func recombineChannels(taskChanIn <-chan *base.Task, taskChanIn2 <-chan *base.Task, taskChanOut chan<- *base.Task, unmatchedTasks *sync.Map) {
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

func tryMatchTask(task *base.Task, taskChanOut chan<- *base.Task, unmatchedTasks *sync.Map) {
	if _, exists := unmatchedTasks.Load(task.ID); exists {
		taskChanOut <- task
		unmatchedTasks.Delete(task.ID)
	} else {
		unmatchedTasks.Store(task.ID, task)
	}
}

func logDoneTasks(taskChanIn <-chan *base.Task) {
	for task := range taskChanIn {
		logger.Debugf("Done processing and updating task: %+v", task.ID)
	}
}
