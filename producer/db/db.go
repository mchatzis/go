package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDBPool() (*pgxpool.Pool, error) {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	return dbpool, nil
}
