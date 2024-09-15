package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"producer/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbpool, err := initDBPool()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	fetchAndPrintTasks(dbpool)

	for {
		fmt.Println("Hello producer")
		time.Sleep(5 * time.Second)
	}
}

func initDBPool() (*pgxpool.Pool, error) {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	return dbpool, nil
}

func fetchAndPrintTasks(dbpool *pgxpool.Pool) {
	queries := db.New(dbpool)
	tasks, err := queries.GetTasks(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(tasks)
}
