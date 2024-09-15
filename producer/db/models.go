// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type TaskState string

const (
	TaskStatePending    TaskState = "pending"
	TaskStateInProgress TaskState = "in_progress"
	TaskStateCompleted  TaskState = "completed"
	TaskStateFailed     TaskState = "failed"
)

func (e *TaskState) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TaskState(s)
	case string:
		*e = TaskState(s)
	default:
		return fmt.Errorf("unsupported scan type for TaskState: %T", src)
	}
	return nil
}

type NullTaskState struct {
	TaskState TaskState
	Valid     bool // Valid is true if TaskState is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTaskState) Scan(value interface{}) error {
	if value == nil {
		ns.TaskState, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TaskState.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTaskState) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TaskState), nil
}

type Task struct {
	ID             int32
	Type           pgtype.Int4
	Value          pgtype.Int4
	State          NullTaskState
	Creationtime   pgtype.Float8
	Lastupdatetime pgtype.Float8
}
