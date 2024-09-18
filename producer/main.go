package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mchatzis/go/producer/db"
)

func main() {
	connection := &db.Connection{}
	if err := connection.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	Run(connection.Pool)

	for {
		fmt.Println("Hello producer")
		time.Sleep(5 * time.Second)
	}
}
