package testing

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchatzis/go/producer/pkg/logging"
)

var logger = logging.GetLogger()

type TestDB struct {
	*pgxpool.Pool
}

var (
	testDbInstance *TestDB
	once           sync.Once
)

func GetDB() *TestDB {
	once.Do(func() {
		ctx := context.Background()
		pool, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
		if err != nil {
			logger.Fatalf("Failed to create connection pool: %v", err)
		}

		testDbInstance = &TestDB{Pool: pool}
	})
	return testDbInstance
}

func (db *TestDB) SetSchema(schema string) {
	_, err := db.Exec(context.Background(), schema)
	if err != nil {
		logger.Fatalf("Failed to create table: %v", err)
	}
}

func (db *TestDB) TearDownDB(dropQuery string) {
	if db != nil && db.Pool != nil {
		_, err := db.Exec(context.Background(), dropQuery)
		if err != nil {
			logger.Fatalf("Failed to tear down DB: %v", err)
		}
		db.Exec(context.Background(), "SELECT * FROM tasks;")
		db.Pool.Close()
	}
}

func (db *TestDB) WithTx(t *testing.T, commit bool, testFunc func(*testing.T, pgx.Tx) error) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			t.Fatalf("Panic in test: %v", r)
		} else if commit {
			err = tx.Commit(ctx)
			if err != nil {
				t.Fatalf("Failed to commit transaction: %v", err)
			}
		} else {
			err = tx.Rollback(ctx)
			if err != nil && err != pgx.ErrTxClosed {
				t.Fatalf("Failed to rollback transaction: %v", err)
			}
		}
	}()

	return testFunc(t, tx)
}
