package main

import (
	"context"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchatzis/go/producer/internal/producer"
	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/mchatzis/go/producer/pkg/monitoring"
	"github.com/mchatzis/go/producer/pkg/sqlc"
)

type Config struct {
	LogLevel string
	DBURL    string
}

var logger = logging.GetLogger()

func main() {
	config := parseFlags()
	setupMonitoring(6060, 2112)
	setupLogging(config.LogLevel)

	err := setupLogging(config.LogLevel)
	if err != nil {
		logger.Fatalf("Failed to set up logging: %v", err)
	}

	dbpool, err := setupDatabase(os.Getenv("DB_URL"))
	if err != nil {
		logger.Fatalf("Failed to set up database: %v", err)
	}
	defer dbpool.Close()
	queries := sqlc.New(dbpool)

	logger.Info("Starting task production...")
	go producer.Produce(queries)
	go monitoring.UpdateMetrics(queries)

	select {}
}

func parseFlags() Config {
	logLevelFlag := flag.String("loglevel", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	return Config{
		LogLevel: *logLevelFlag,
	}
}

func setupMonitoring(metricsPort, profilingPort int) {
	monitoring.RegisterCollectors()
	go monitoring.ExposeMetrics(metricsPort)
	go monitoring.ExposeProfiling(profilingPort)
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
