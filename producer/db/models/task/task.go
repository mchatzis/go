package task

import (
	"context"
	"log"
	"time"

	"producer/db"
	"producer/db/sqlc"
)

// Task database model
type Task struct {
	ID             int       `json:"id"`               // Task ID (Auto-incremented by the database)
	Type           int       `json:"type"`             // Task Type (0-9)
	Value          int       `json:"value"`            // Task Value (0-99)
	State          TaskState `json:"state"`            // Current task state (pending, in_progress, completed, failed)
	CreationTime   time.Time `json:"creation_time"`    // Time of creation (Go uses time.Time for precise time representation)
	LastUpdateTime time.Time `json:"last_update_time"` // Last updated time
}

func FetchAndPrintTasks(database *db.Database) {
	queries := sqlc.New(database.Pool)
	tasks, err := queries.GetTasks(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(tasks)
}
