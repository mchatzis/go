package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func (db *Database) Init() error {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}
	db.Pool = dbpool
	return nil
}

func (db *Database) Close() {
	db.Pool.Close()
}
