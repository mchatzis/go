// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"
)

const getTasks = `-- name: GetTasks :many
SELECT id, type, value, state, creationtime, lastupdatetime FROM tasks
`

func (q *Queries) GetTasks(ctx context.Context) ([]Task, error) {
	rows, err := q.db.Query(ctx, getTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Value,
			&i.State,
			&i.Creationtime,
			&i.Lastupdatetime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
