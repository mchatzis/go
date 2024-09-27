package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchatzis/go/producer/internal/producer"
	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/mchatzis/go/producer/pkg/monitoring"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

var Version string

type Config struct {
	LogLevel   string
	DBURL      string
	MaxBacklog int
}

var logger = logging.GetLogger()

func main() {
	config := parseFlags()

	setupLogging(config.LogLevel)
	err := setupLogging(config.LogLevel)
	if err != nil {
		logger.Fatalf("Failed to set up logging: %v", err)
	}

	monitoring.RegisterCollectors()
	go monitoring.StartServer(os.Getenv("PRODUCER_METRICS_PORT"))

	dbpool, err := setupDatabase(os.Getenv("DB_URL"))
	if err != nil {
		logger.Fatalf("Failed to set up database: %v", err)
	}
	defer dbpool.Close()
	queries := sqlc.New(dbpool)

	logger.Info("Starting task production...")
	go producer.Produce(queries, config.MaxBacklog)
	go monitoring.UpdateMetrics(queries)

	select {}
}

func parseFlags() Config {
	logLevelFlag := flag.String("loglevel", "info", "Log level (debug, info, warn, error)")
	maxBacklogFlag := flag.Int("max_backlog", 50, "Maximum back log, after that stop sending to consumer")
	flag.Parse()

	return Config{
		LogLevel:   *logLevelFlag,
		MaxBacklog: *maxBacklogFlag,
	}
}

func setupLogging(logLevelFlag string) error {
	level, err := getLogLevel(logLevelFlag)
	if err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}
	logging.SetLogLevel(level)
	return nil
}

func setupDatabase(dbURL string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create db pool: %v", err)
	}
	return dbpool, nil
}

func getLogLevel(logLevelFlag string) (logging.LogLevel, error) {
	switch logLevelFlag {
	case "debug":
		return logging.DEBUG, nil
	case "info":
		return logging.INFO, nil
	case "warn":
		return logging.WARN, nil
	case "error":
		return logging.ERROR, nil
	default:
		return logging.LogLevel(0), fmt.Errorf("%s is not a valid loglevel flag value", logLevelFlag)
	}
}
