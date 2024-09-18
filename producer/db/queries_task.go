package db

import (
	"context"

	"github.com/mchatzis/go/producer/db/sqlc"
)

func CreateTask(queries *sqlc.Queries, task sqlc.Task) error {
	params := sqlc.CreateTaskParams(task)
	err := queries.CreateTask(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTasks(queries *sqlc.Queries) ([]sqlc.Task, error) {
	tasks, err := queries.GetTasks(context.Background())
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
