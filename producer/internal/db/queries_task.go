package db

import (
	"context"

	"github.com/mchatzis/go/producer/pkg/sqlc"
)

func CreateTask(queries *sqlc.Queries, task sqlc.Task) error {
	params := sqlc.CreateTaskParams(task)
	err := queries.CreateTask(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}
