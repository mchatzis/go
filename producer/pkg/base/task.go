package base

import (
	"github.com/mchatzis/go/producer/pkg/grpc"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

type TaskState string

const (
	TaskStatePending    TaskState = "pending"
	TaskStateProcessing TaskState = "processing"
	TaskStateDone       TaskState = "done"
	TaskStateFailed     TaskState = "failed"
)

type Task struct {
	ID             int32
	Type           int32
	Value          int32
	State          TaskState
	Creationtime   float64
	Lastupdatetime float64
	FutureField1   string
	FutureField2   string
}

func (t *Task) FromSQLCTask(sqlcTask *sqlc.Task) {
	t.ID = sqlcTask.ID
	t.Type = sqlcTask.Type
	t.Value = sqlcTask.Value
	t.State = mapFromSqlcState(sqlcTask.State)
	t.Creationtime = sqlcTask.Creationtime
	t.Lastupdatetime = sqlcTask.Lastupdatetime
}

func (t *Task) ToSQLCTask() *sqlc.Task {
	return &sqlc.Task{
		ID:             t.ID,
		Type:           t.Type,
		Value:          t.Value,
		State:          mapToSqlcState(t.State),
		Creationtime:   t.Creationtime,
		Lastupdatetime: t.Lastupdatetime,
	}
}

func (t *Task) FromGRPCTask(GrpcTask *grpc.Task) {
	t.ID = GrpcTask.Id
	t.Type = GrpcTask.Type
	t.Value = GrpcTask.Value
	t.State = mapFromGrpcState(GrpcTask.State)
	t.Creationtime = GrpcTask.CreationTime
	t.Lastupdatetime = GrpcTask.LastUpdateTime
}

func (t *Task) ToGRPCTask() *grpc.Task {
	return &grpc.Task{
		Id:             t.ID,
		Type:           t.Type,
		Value:          t.Value,
		State:          mapToGrpcState(t.State),
		CreationTime:   t.Creationtime,
		LastUpdateTime: t.Lastupdatetime,
	}
}

func mapFromGrpcState(state grpc.TaskState) TaskState {
	switch state {
	case grpc.TaskState_PENDING:
		return TaskStatePending
	case grpc.TaskState_PROCESSING:
		return TaskStateProcessing
	case grpc.TaskState_DONE:
		return TaskStateDone
	case grpc.TaskState_FAILED:
		return TaskStateFailed
	default:
		panic("Unexpected grpc.TaskState value")
	}
}

func mapToGrpcState(state TaskState) grpc.TaskState {
	switch state {
	case TaskStatePending:
		return grpc.TaskState_PENDING
	case TaskStateProcessing:
		return grpc.TaskState_PROCESSING
	case TaskStateDone:
		return grpc.TaskState_DONE
	case TaskStateFailed:
		return grpc.TaskState_FAILED
	default:
		panic("Unexpected base.TaskState value")
	}
}

func mapFromSqlcState(state sqlc.TaskState) TaskState {
	switch state {
	case sqlc.TaskStatePending:
		return TaskStatePending
	case sqlc.TaskStateProcessing:
		return TaskStateProcessing
	case sqlc.TaskStateDone:
		return TaskStateDone
	case sqlc.TaskStateFailed:
		return TaskStateFailed
	default:
		panic("Unexpected sqlc.TaskState value")
	}
}

func mapToSqlcState(state TaskState) sqlc.TaskState {
	switch state {
	case TaskStatePending:
		return sqlc.TaskStatePending
	case TaskStateProcessing:
		return sqlc.TaskStateProcessing
	case TaskStateDone:
		return sqlc.TaskStateDone
	case TaskStateFailed:
		return sqlc.TaskStateFailed
	default:
		panic("Unexpected base.TaskState value")
	}
}
