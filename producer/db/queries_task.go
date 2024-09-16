package db

import (
	"context"
	"log"
	"producer/db/sqlc"
)

func FetchAndPrintTasks(database *Database) {
	queries := sqlc.New(database.Pool)
	tasks, err := queries.GetTasks(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(tasks[0].Value.Valid)
	log.Println(tasks)
}
