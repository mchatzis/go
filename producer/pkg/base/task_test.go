package base

import (
	"testing"

	"github.com/mchatzis/go/producer/pkg/grpc"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

func TestFromSQLCTask(t *testing.T) {
	sqlcTask := &sqlc.Task{
		ID:             1,
		Type:           2,
		Value:          3,
		State:          sqlc.TaskStateDone,
		Creationtime:   1620000000,
		Lastupdatetime: 1620000001,
	}

	task := &Task{}
	task.FromSQLCTask(sqlcTask)

	if task.ID != sqlcTask.ID || task.Type != sqlcTask.Type || task.Value != sqlcTask.Value ||
		task.State != TaskStateDone || task.Creationtime != sqlcTask.Creationtime ||
		task.Lastupdatetime != sqlcTask.Lastupdatetime {
		t.Errorf("FromSQLCTask failed: got %+v, want %+v", task, sqlcTask)
	}
}

func TestToSQLCTask(t *testing.T) {
	task := &Task{
		ID:             1,
		Type:           2,
		Value:          3,
		State:          TaskStateDone,
		Creationtime:   1620000000,
		Lastupdatetime: 1620000001,
	}

	sqlcTask := task.ToSQLCTask()

	if sqlcTask.ID != task.ID || sqlcTask.Type != task.Type || sqlcTask.Value != task.Value ||
		sqlcTask.State != sqlc.TaskStateDone || sqlcTask.Creationtime != task.Creationtime ||
		sqlcTask.Lastupdatetime != task.Lastupdatetime {
		t.Errorf("ToSQLCTask failed: got %+v, want %+v", sqlcTask, task)
	}
}

func TestFromGRPCTask(t *testing.T) {
	grpcTask := &grpc.Task{
		Id:             1,
		Type:           2,
		Value:          3,
		State:          grpc.TaskState_DONE,
		CreationTime:   1620000000,
		LastUpdateTime: 1620000001,
	}

	task := &Task{}
	task.FromGRPCTask(grpcTask)

	if task.ID != grpcTask.Id || task.Type != grpcTask.Type || task.Value != grpcTask.Value ||
		task.State != TaskStateDone || task.Creationtime != grpcTask.CreationTime ||
		task.Lastupdatetime != grpcTask.LastUpdateTime {
		t.Errorf("FromGRPCTask failed: got %+v, want %+v", task, grpcTask)
	}
}

func TestToGRPCTask(t *testing.T) {
	task := &Task{
		ID:             1,
		Type:           2,
		Value:          3,
		State:          TaskStateDone,
		Creationtime:   1620000000,
		Lastupdatetime: 1620000001,
	}

	grpcTask := task.ToGRPCTask()

	if grpcTask.Id != task.ID || grpcTask.Type != task.Type || grpcTask.Value != task.Value ||
		grpcTask.State != grpc.TaskState_DONE || grpcTask.CreationTime != task.Creationtime ||
		grpcTask.LastUpdateTime != task.Lastupdatetime {
		t.Errorf("ToGRPCTask failed: got %+v, want %+v", grpcTask, task)
	}
}

func TestMapFromGrpcState(t *testing.T) {
	tests := []struct {
		input    grpc.TaskState
		expected TaskState
	}{
		{grpc.TaskState_PENDING, TaskStatePending},
		{grpc.TaskState_PROCESSING, TaskStateProcessing},
		{grpc.TaskState_DONE, TaskStateDone},
		{grpc.TaskState_FAILED, TaskStateFailed},
	}

	for _, test := range tests {
		result := mapFromGrpcState(test.input)
		if result != test.expected {
			t.Errorf("mapFromGrpcState(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestMapToGrpcState(t *testing.T) {
	tests := []struct {
		input    TaskState
		expected grpc.TaskState
	}{
		{TaskStatePending, grpc.TaskState_PENDING},
		{TaskStateProcessing, grpc.TaskState_PROCESSING},
		{TaskStateDone, grpc.TaskState_DONE},
		{TaskStateFailed, grpc.TaskState_FAILED},
	}

	for _, test := range tests {
		result := mapToGrpcState(test.input)
		if result != test.expected {
			t.Errorf("mapToGrpcState(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestMapFromSqlcState(t *testing.T) {
	tests := []struct {
		input    sqlc.TaskState
		expected TaskState
	}{
		{sqlc.TaskStatePending, TaskStatePending},
		{sqlc.TaskStateProcessing, TaskStateProcessing},
		{sqlc.TaskStateDone, TaskStateDone},
		{sqlc.TaskStateFailed, TaskStateFailed},
	}

	for _, test := range tests {
		result := mapFromSqlcState(test.input)
		if result != test.expected {
			t.Errorf("mapFromSqlcState(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestMapToSqlcState(t *testing.T) {
	tests := []struct {
		input    TaskState
		expected sqlc.TaskState
	}{
		{TaskStatePending, sqlc.TaskStatePending},
		{TaskStateProcessing, sqlc.TaskStateProcessing},
		{TaskStateDone, sqlc.TaskStateDone},
		{TaskStateFailed, sqlc.TaskStateFailed},
	}

	for _, test := range tests {
		result := mapToSqlcState(test.input)
		if result != test.expected {
			t.Errorf("mapToSqlcState(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}
