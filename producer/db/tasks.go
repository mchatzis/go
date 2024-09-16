package db

import (
	"context"
	"log"
	db "producer/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

func FetchAndPrintTasks(dbpool *pgxpool.Pool) {
	queries := db.New(dbpool)
	tasks, err := queries.GetTasks(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(tasks)
}
