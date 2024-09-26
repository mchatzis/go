package consumer

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mchatzis/go/producer/pkg/sqlc"
	"github.com/mchatzis/go/producer/pkg/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateTasksStateInDb(t *testing.T) {
	t.Run("Successful update", func(t *testing.T) {
		mockDB := new(testing.MockDBTX)
		queries := sqlc.New(mockDB)

		expectedState := sqlc.TaskStateProcessing
		expectedTaskID := int32(1)

		mockDB.On("Exec", mock.Anything, mock.Anything, []interface{}{expectedState, expectedTaskID}).
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
		mockDB := new(testing.MockDBTX)
		queries := sqlc.New(mockDB)

		expectedState := sqlc.TaskStateProcessing
		expectedTaskID := int32(1)
		expectedError := errors.New("Failed to perform database operation update state by id")

		mockDB.On("Exec", mock.Anything, mock.Anything, []interface{}{expectedState, expectedTaskID}).
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
