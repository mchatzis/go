package main

import (
	"fmt"
	"os"
	"time"

	"producer/db"
)

func main() {
	dbpool, err := db.InitDBPool()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	db.FetchAndPrintTasks(dbpool)

	for {
		fmt.Println("Hello producer")
		time.Sleep(5 * time.Second)
	}
}
