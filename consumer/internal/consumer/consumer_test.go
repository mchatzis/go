package consumer

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mchatzis/go/producer/pkg/sqlc"
	prod_testing "github.com/mchatzis/go/producer/pkg/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateTasksStateInDb(t *testing.T) {
	t.Run("Successful update", func(t *testing.T) {
		mockDB := new(prod_testing.MockDBTX)
		queries := sqlc.New(mockDB)

		expectedState := sqlc.TaskStateProcessing
		expectedTaskID := int32(1)

		mockDB.On("Exec", mock.Anything, mock.Anything, expectedState, expectedTaskID).
			Return(pgconn.CommandTag{}, nil)

		taskChanIn := make(chan *sqlc.Task)
		taskChanOut := make(chan *sqlc.Task)

		go updateTasksStateInDb(taskChanIn, taskChanOut, expectedState, queries)

		inputTask := &sqlc.Task{ID: expectedTaskID, State: sqlc.TaskStatePending}
		taskChanIn <- inputTask
		close(taskChanIn)

		outputTask := <-taskChanOut

		mockDB.AssertExpectations(t)
		assert.Equal(t, inputTask.ID, outputTask.ID, "Output task ID should be identical to input task")
		assert.Equal(t, inputTask.State, outputTask.State, "Output task State should be identical to input task")
	})

	t.Run("Failed update", func(t *testing.T) {
		mockDB := new(prod_testing.MockDBTX)
		queries := sqlc.New(mockDB)

		expectedState := sqlc.TaskStateProcessing
		expectedTaskID := int32(1)
		expectedError := errors.New("Failed to perform database operation update state by id")

		mockDB.On("Exec", mock.Anything, mock.Anything, expectedState, expectedTaskID).
			Return(pgconn.CommandTag{}, expectedError)

		taskChanIn := make(chan *sqlc.Task)
		taskChanOut := make(chan *sqlc.Task)

		go updateTasksStateInDb(taskChanIn, taskChanOut, expectedState, queries)

		inputTask := &sqlc.Task{ID: expectedTaskID, State: sqlc.TaskStatePending}
		taskChanIn <- inputTask
		close(taskChanIn)

		outputTask := <-taskChanOut

		mockDB.AssertExpectations(t)
		assert.Equal(t, inputTask, outputTask, "Output task should be identical to input task even on error")
	})
}

func TestDistributeIncomingTasks(t *testing.T) {
	testCases := []struct {
		name     string
		inputLen int
	}{
		{"Single task", 1},
		{"Multiple tasks", 5},
		{"No tasks", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskChanIn := make(chan *sqlc.Task, tc.inputLen)
			taskChanOut := make(chan *sqlc.Task, tc.inputLen)
			taskChanOut2 := make(chan *sqlc.Task, tc.inputLen)

			for i := 0; i < tc.inputLen; i++ {
				taskChanIn <- &sqlc.Task{
					ID:    int32(i),
					State: sqlc.TaskStatePending,
				}
			}
			close(taskChanIn)

			var rateLimitDuration = 0 * time.Millisecond
			go distributeIncomingTasksWithRateLimit(taskChanIn, taskChanOut, taskChanOut2, rateLimitDuration)

			for i := 0; i < tc.inputLen; i++ {
				select {
				case task := <-taskChanOut:
					if task.State != sqlc.TaskStateProcessing {
						t.Errorf("Expected task state to be %s, got %s", sqlc.TaskStateProcessing, task.State)
					}
				case <-time.After(time.Second):
					t.Error("Timeout waiting for task on taskChanOut")
				}

				select {
				case task := <-taskChanOut2:
					if task.State != sqlc.TaskStateProcessing {
						t.Errorf("Expected task state to be %s, got %s", sqlc.TaskStateProcessing, task.State)
					}
				case <-time.After(time.Second):
					t.Error("Timeout waiting for task on taskChanOut2")
				}
			}

			select {
			case <-taskChanOut:
				t.Error("Unexpected extra task on taskChanOut")
			case <-taskChanOut2:
				t.Error("Unexpected extra task on taskChanOut2")
			default:
			}
		})
	}
}

func TestDistributeIncomingTasksRateLimit(t *testing.T) {
	const numTasks = 5
	taskChanIn := make(chan *sqlc.Task, numTasks)
	taskChanOut := make(chan *sqlc.Task, numTasks)
	taskChanOut2 := make(chan *sqlc.Task, numTasks)

	for i := 1; i <= numTasks; i++ {
		taskChanIn <- &sqlc.Task{ID: int32(i), State: sqlc.TaskStatePending}
	}
	close(taskChanIn)

	var rateLimitDuration = 100 * time.Millisecond
	start := time.Now()
	go distributeIncomingTasksWithRateLimit(taskChanIn, taskChanOut, taskChanOut2, rateLimitDuration)

	receivedTasks := 0
	for receivedTasks < 2*numTasks {
		select {
		case <-taskChanOut:
			receivedTasks++
		case <-taskChanOut2:
			receivedTasks++
		}
	}

	duration := time.Since(start)

	expectedDuration := time.Duration(numTasks-1) * rateLimitDuration
	if duration < expectedDuration {
		t.Errorf("Expected rate limiting of at least %v, got %v", expectedDuration, duration)
	}
}

type MockProcessor struct {
	mock.Mock
}

func (m *MockProcessor) Process(task *sqlc.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func TestProcessTasks(t *testing.T) {
	mockProc := new(MockProcessor)

	mockProc.On("Process", mock.AnythingOfType("*sqlc.Task")).Return(nil).Once()
	mockProc.On("Process", mock.AnythingOfType("*sqlc.Task")).Return(errors.New("processing error")).Once()
	mockProc.On("Process", mock.AnythingOfType("*sqlc.Task")).Return(nil).Once()

	taskChanIn := make(chan *sqlc.Task, 3)
	taskChanOut := make(chan *sqlc.Task, 3)

	go processTasks(taskChanIn, taskChanOut, mockProc.Process)

	tasks := []*sqlc.Task{
		{ID: 1, State: sqlc.TaskStatePending},
		{ID: 2, State: sqlc.TaskStatePending},
		{ID: 3, State: sqlc.TaskStatePending},
	}

	for _, task := range tasks {
		taskChanIn <- task
	}
	close(taskChanIn)

	processedTasks := make(map[int32]*sqlc.Task)
	for i := 0; i < 2; i++ { // We expect 2 successful tasks
		select {
		case task := <-taskChanOut:
			processedTasks[task.ID] = task
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for processed task")
		}
	}

	assert.Contains(t, processedTasks, int32(1), "Task 1 should be in the processed tasks")
	assert.Equal(t, sqlc.TaskStateDone, processedTasks[1].State, "Task 1 should be done")

	assert.NotContains(t, processedTasks, int32(2), "Task 2 should not be in the processed tasks")
	assert.Equal(t, sqlc.TaskStateFailed, tasks[1].State, "Task 2 should be marked as failed")

	assert.Contains(t, processedTasks, int32(3), "Task 3 should be in the processed tasks")
	assert.Equal(t, sqlc.TaskStateDone, processedTasks[3].State, "Task 3 should be done")

	mockProc.AssertExpectations(t)
}

func TestPretendToProcess(t *testing.T) {
	task := &sqlc.Task{Value: 100}

	start := time.Now()
	err := pretendToProcess(task)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, duration.Milliseconds(), int64(100), "Expected at least 100ms of processing time")
}
