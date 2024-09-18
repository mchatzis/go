package main

import (
	"log"
	"time"

	"github.com/mchatzis/go/producer/db/sqlc"
)

func main() {
	taskChan := make(chan *sqlc.Task, 100)
	go ListenForTasks(taskChan)

	for i := 0; i < 15; i++ {
		go processTask(taskChan)
	}

	select {}
}

func processTask(taskChan chan *sqlc.Task) {
	for task := range taskChan {
		log.Printf("Processing task: %+v", *task)
		time.Sleep(time.Duration(task.Value) * time.Millisecond)
		log.Printf("Finished processing task: %+v", *task)
	}
}
