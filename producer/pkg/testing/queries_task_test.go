package testing

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

func TestMain(m *testing.M) {
	db := GetDB()
	db.SetSchema(schema)
	code := m.Run()
	db.TearDownDB(dropQuery)
	os.Exit(code)
}

func TestSchema(t *testing.T) {
	db := GetDB()

	now := float64(time.Now().UnixNano()) / 1e9

	testCases := []struct {
		name           string
		id             int
		taskType       int
		value          int
		state          string
		creationTime   float64
		lastUpdateTime float64
		expectError    bool
		commit         bool
	}{
		{"Valid insertion", 1, 5, 50, "pending", now, now, false, true},
		{"Zero ID", 0, 5, 50, "pending", now, now, true, false},
		{"Negative ID", -1, 5, 50, "pending", now, now, true, false},
		{"Type too low", 2, -1, 50, "pending", now, now, true, false},
		{"Type too high", 3, 10, 50, "pending", now, now, true, false},
		{"Value too low", 4, 5, -1, "pending", now, now, true, false},
		{"Value too high", 5, 5, 100, "pending", now, now, true, false},
		{"Invalid state", 6, 5, 50, "invalid", now, now, true, false},
		{"Duplicate ID", 1, 5, 50, "pending", now, now, true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := db.WithTx(t, tc.commit, func(t *testing.T, tx pgx.Tx) error {
				txQueries := sqlc.New(tx)
				task := sqlc.Task{
					ID:             int32(tc.id),
					Type:           int32(tc.taskType),
					Value:          int32(tc.value),
					State:          sqlc.TaskState(tc.state),
					Creationtime:   tc.creationTime,
					Lastupdatetime: tc.lastUpdateTime,
				}
				params := sqlc.CreateTaskParams(task)
				return txQueries.CreateTask(context.Background(), params)
			})
			if tc.expectError && err == nil {
				t.Errorf("Expected error for %s, but got none", tc.name)
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.name, err)
			}
		})
	}
}

const schema = `DO $$
BEGIN
	CREATE TYPE task_state AS ENUM ('pending', 'processing', 'done', 'failed');
END $$;

DO $$
BEGIN
	CREATE TABLE tasks (
		ID INT PRIMARY KEY CHECK (ID > 0),
		Type INT CHECK (Type BETWEEN 0 AND 9) NOT NULL,
		Value INT CHECK (Value BETWEEN 0 AND 99) NOT NULL,
		State task_state NOT NULL,
		CreationTime FLOAT NOT NULL,
		LastUpdateTime FLOAT NOT NULL
	);
END $$;`

const dropQuery = `DO $$
BEGIN
	DROP TABLE IF EXISTS tasks;
END $$;

DO $$
BEGIN
	DROP TYPE IF EXISTS task_state;
END $$;
`
