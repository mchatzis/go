package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mchatzis/go/producer/internal/db"
	"github.com/mchatzis/go/producer/internal/producer"
)

func main() {
	connection := &db.Connection{}
	if err := connection.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	producer.Produce(connection.Pool)

	for {
		time.Sleep(5 * time.Second)
	}
}
