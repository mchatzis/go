package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	Pool *pgxpool.Pool
}

func (db *Connection) Init() error {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}
	db.Pool = dbpool
	return nil
}

func (db *Connection) Close() {
	db.Pool.Close()
}
