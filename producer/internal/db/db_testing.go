package db

import (
	"database/sql"
	"sync"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mchatzis/go/producer/pkg/logging"
)

var logger = logging.GetLogger()

type TestDB struct {
	*sql.DB
}

var (
	testDbInstance *TestDB
	once           sync.Once
)

func GetDB() *TestDB {
	once.Do(func() {
		var err error
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			logger.Fatalf("Failed to open database: %v", err)
		}

		testDbInstance = &TestDB{DB: db}

		schemaQuery := `
			CREATE TABLE tasks (
				ID INTEGER UNIQUE NOT NULL CHECK (ID > 0),
				Type INTEGER CHECK (Type BETWEEN 0 AND 9) NOT NULL,
				Value INTEGER CHECK (Value BETWEEN 0 AND 99) NOT NULL,
				State TEXT CHECK (State IN ('pending', 'processing', 'done', 'failed')) NOT NULL,
				CreationTime REAL NOT NULL,
				LastUpdateTime REAL NOT NULL
			);
		`

		_, err = db.Exec(schemaQuery)
		if err != nil {
			logger.Fatalf("Failed to create table: %v", err)
		}
	})
	return testDbInstance
}

func (db *TestDB) TearDownDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Fatalf("Failed to close database: %v", err)
		}
		db = nil
	}
}

func (db *TestDB) WithTx(t *testing.T, commit bool, testFunc func(*testing.T, *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			t.Fatalf("Panic in test: %v", r)
		} else if commit {
			err = tx.Commit()
			if err != nil {
				t.Fatalf("Failed to commit transaction: %v", err)
			}
		} else {
			err = tx.Rollback()
			if err != nil && err != sql.ErrTxDone {
				t.Fatalf("Failed to rollback transaction: %v", err)
			}
		}
	}()

	return testFunc(t, tx)
}
