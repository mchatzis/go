package main

import (
	"context"
	"flag"
	_ "net/http/pprof"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchatzis/go/producer/internal/producer"
	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/mchatzis/go/producer/pkg/monitoring"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

var logger = logging.GetLogger()

func init() {
	monitoring.RegisterCollectors()
	go monitoring.ExposeMetrics(2112)
	go monitoring.ExposeProfiling(6060)
}

func main() {
	logLevelFlag := flag.String("loglevel", "info", "Log level (debug, info, warn, error)")
	flag.Parse()
	setLogerLevel(*logLevelFlag)

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		logger.Fatalf("Failed to create db pool with error: %v\n", err)
	}
	defer dbpool.Close()
	queries := sqlc.New(dbpool)

	logger.Info("Starting task production...")
	go producer.Produce(queries)

	go monitoring.UpdateMetrics(queries)

	select {}
}

func setLogerLevel(logLevelFlag string) {
	switch logLevelFlag {
	case "debug":
		logging.SetLogLevel(logging.DEBUG)
	case "info":
		logging.SetLogLevel(logging.INFO)
	case "warn":
		logging.SetLogLevel(logging.WARN)
	case "error":
		logging.SetLogLevel(logging.ERROR)
	default:
		logger.Fatalf("%+v is not a valid loglevel flag value", logLevelFlag)
	}
}
