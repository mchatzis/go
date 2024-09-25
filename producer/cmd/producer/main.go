package main

import (
	"context"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchatzis/go/producer/internal/producer"
	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/mchatzis/go/producer/pkg/monitoring"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

var logger = logging.GetLogger()

func init() {
	go func() {
		logger.Error(http.ListenAndServe(":6060", nil))
	}()
	monitoring.RegisterCollectors()
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

	//Wait a sec for postgres container to finish spin up, should use a docker healthcheck instead in the future
	time.Sleep(time.Second)

	logger.Info("Starting task production...")
	go producer.Produce(queries)

	go monitoring.ExposeMetrics()
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
