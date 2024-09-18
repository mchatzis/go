package main

import (
	"log"
)

func main() {
	taskChan := make(chan *Task, 100)
	go ListenForTasks(taskChan)

	for task := range taskChan {
		log.Printf("Processing task: %+v", *task)
	}
}
