package main

import (
	"fmt"
	"os"
	"time"

	"producer/db"
)

func main() {
	database := &db.Database{}
	if err := database.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	db.FetchAndPrintTasks(database)

	for {
		fmt.Println("Hello producer")
		time.Sleep(5 * time.Second)
	}
}
