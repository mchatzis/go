package main

import (
	"fmt"
	"os"

	db "github.com/mchatzis/go/producer/db"
	"github.com/mchatzis/go/producer/pkg/sqlc"
	"github.com/mchatzis/go/consumer/internal/consumer"
)

func main() {
	connection := &db.Connection{}
	if err := connection.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer connection.Close()
	queries := sqlc.New(connection.Pool)

	go consumer.consume(queries)

	select {}
}
