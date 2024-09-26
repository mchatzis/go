package main

import (
	"context"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchatzis/go/consumer/internal/consumer"
	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/mchatzis/go/producer/pkg/monitoring"
	"github.com/mchatzis/go/producer/pkg/sqlc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var logger = logging.GetLogger()

type Config struct {
	LogLevel string
	DBURL    string
}

func main() {
	config := parseFlags()
	setupMonitoring(6061, 2112)
	setupLogging(config.LogLevel)

	dbpool := setupDatabase(os.Getenv("DB_URL"))
	defer dbpool.Close()
	queries := sqlc.New(dbpool)

	go consumer.HandleIncomingTasks(queries)
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
	promauto.NewGauge(prometheus.GaugeOpts{
		Name: "service_up",
		Help: "Indicates whether the service is up (1) or down (0)",
	}).Set(1)
	go monitoring.ExposeMetrics(metricsPort)
	go monitoring.ExposeProfiling(profilingPort)
}

func setupLogging(logLevelFlag string) {
	logLevel, err := getLogLevel(logLevelFlag)
	if err != nil {
		logger.Fatalf("Invalid log level: %v", err)
	}
	logging.SetLogLevel(logLevel)
}

func setupDatabase(dbURL string) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatalf("Failed to create db pool with error: %v\n", err)
	}
	return dbpool
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
