// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package sqlc

import (
	"context"
)

const countOfTasksInState = `-- name: CountOfTasksInState :one
SELECT COUNT(*) FROM tasks WHERE state = $1
`

func (q *Queries) CountOfTasksInState(ctx context.Context, state TaskState) (int64, error) {
	row := q.db.QueryRow(ctx, countOfTasksInState, state)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createTask = `-- name: CreateTask :exec
INSERT INTO tasks (id, type, value, state, creationtime, lastupdatetime) VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateTaskParams struct {
	ID             int32
	Type           int32
	Value          int32
	State          TaskState
	Creationtime   float64
	Lastupdatetime float64
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) error {
	_, err := q.db.Exec(ctx, createTask,
		arg.ID,
		arg.Type,
		arg.Value,
		arg.State,
		arg.Creationtime,
		arg.Lastupdatetime,
	)
	return err
}

const getCountOfDoneTasksByType = `-- name: GetCountOfDoneTasksByType :many
SELECT type, COUNT(*) FROM tasks WHERE state='done' GROUP BY type ORDER BY type
`

type GetCountOfDoneTasksByTypeRow struct {
	Type  int32
	Count int64
}

func (q *Queries) GetCountOfDoneTasksByType(ctx context.Context) ([]GetCountOfDoneTasksByTypeRow, error) {
	rows, err := q.db.Query(ctx, getCountOfDoneTasksByType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCountOfDoneTasksByTypeRow
	for rows.Next() {
		var i GetCountOfDoneTasksByTypeRow
		if err := rows.Scan(&i.Type, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCountOfTasksByState = `-- name: GetCountOfTasksByState :many
SELECT state, COUNT(*) FROM tasks GROUP BY state
`

type GetCountOfTasksByStateRow struct {
	State TaskState
	Count int64
}

func (q *Queries) GetCountOfTasksByState(ctx context.Context) ([]GetCountOfTasksByStateRow, error) {
	rows, err := q.db.Query(ctx, getCountOfTasksByState)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCountOfTasksByStateRow
	for rows.Next() {
		var i GetCountOfTasksByStateRow
		if err := rows.Scan(&i.State, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTotalValueOfDoneTasksByType = `-- name: GetTotalValueOfDoneTasksByType :many
SELECT type, SUM(value) FROM tasks WHERE state='done' GROUP BY type ORDER BY type
`

type GetTotalValueOfDoneTasksByTypeRow struct {
	Type int32
	Sum  int64
}

func (q *Queries) GetTotalValueOfDoneTasksByType(ctx context.Context) ([]GetTotalValueOfDoneTasksByTypeRow, error) {
	rows, err := q.db.Query(ctx, getTotalValueOfDoneTasksByType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTotalValueOfDoneTasksByTypeRow
	for rows.Next() {
		var i GetTotalValueOfDoneTasksByTypeRow
		if err := rows.Scan(&i.Type, &i.Sum); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTaskState = `-- name: UpdateTaskState :exec
UPDATE tasks SET state = $1 WHERE id = $2
`

type UpdateTaskStateParams struct {
	State TaskState
	ID    int32
}

func (q *Queries) UpdateTaskState(ctx context.Context, arg UpdateTaskStateParams) error {
	_, err := q.db.Exec(ctx, updateTaskState, arg.State, arg.ID)
	return err
}
