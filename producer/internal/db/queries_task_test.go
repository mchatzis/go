package db

import (
	"database/sql"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	db := GetDB()
	code := m.Run()
	db.TearDownDB()
	os.Exit(code)
}

func TestCreateTask(t *testing.T) {
	db := GetDB()

	now := float64(time.Now().UnixNano()) / 1e9

	insertTask := func(tx *sql.Tx, id, taskType, value int, state string, creationTime, lastUpdateTime float64) error {
		_, err := tx.Exec("INSERT INTO tasks (ID, Type, Value, State, CreationTime, LastUpdateTime) VALUES (?, ?, ?, ?, ?, ?)",
			id, taskType, value, state, creationTime, lastUpdateTime)
		return err
	}

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
			err := db.WithTx(t, tc.commit, func(t *testing.T, tx *sql.Tx) error {
				return insertTask(tx, tc.id, tc.taskType, tc.value, tc.state, tc.creationTime, tc.lastUpdateTime)
			})
			if tc.expectError && err == nil {
				t.Errorf("Expected error for %s, but got none", tc.name)
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.name, err)
			}
		})
	}

}
