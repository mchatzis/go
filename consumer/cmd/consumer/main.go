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

var Version string

var logger = logging.GetLogger()

type Config struct {
	ShowVersion         bool
	LogLevel            string
	DBURL               string
	RateLimitMultiplier int
}

func main() {
	config := parseFlags()

	if config.ShowVersion {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	err := setupLogging(config.LogLevel)
	if err != nil {
		logger.Fatalf("Failed to set up logging: %v", err)
	}

	promauto.NewGauge(prometheus.GaugeOpts{
		Name: "service_up",
		Help: "Indicates whether the service is up (1) or down (0)",
	}).Set(1)
	go monitoring.ExposeMetrics(os.Getenv("CONSUMER_METRICS_PORT"))

	dbpool, err := setupDatabase(os.Getenv("DB_URL"))
	if err != nil {
		logger.Fatalf("Failed to set up database: %v", err)
	}
	defer dbpool.Close()
	queries := sqlc.New(dbpool)

	go consumer.HandleIncomingTasks(queries, config.RateLimitMultiplier)
	select {}
}

func parseFlags() Config {
	versionFlag := flag.Bool("version", false, "Show version information")
	logLevelFlag := flag.String("loglevel", "info", "Log level (debug, info, warn, error)")
	rateLimitMultiplierFlag := flag.Int("rate_limit_multiplier", 500, "Limit the rate of incoming tasks. Introduces delay of that many milliseconds")
	flag.Parse()

	return Config{
		ShowVersion:         *versionFlag,
		LogLevel:            *logLevelFlag,
		RateLimitMultiplier: *rateLimitMultiplierFlag,
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
