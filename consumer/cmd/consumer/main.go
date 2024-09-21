package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchatzis/go/consumer/internal/consumer"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	queries := sqlc.New(dbpool)

	go consumer.Consume(queries)

	select {}
}
