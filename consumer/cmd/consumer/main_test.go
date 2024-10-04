package main

import (
	"flag"
	"os"
	"testing"

	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestGetLogLevel(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedLevel logging.LogLevel
		expectError   bool
	}{
		{"Debug level", "debug", logging.DEBUG, false},
		{"Info level", "info", logging.INFO, false},
		{"Warn level", "warn", logging.WARN, false},
		{"Error level", "error", logging.ERROR, false},
		{"Invalid level", "invalid", logging.LogLevel(0), true},
		{"Empty string", "", logging.LogLevel(0), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			level, err := getLogLevel(tc.input)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedLevel, level)
			}
		})
	}
}

func TestParseFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	testCases := []struct {
		name           string
		args           []string
		expectedConfig Config
	}{
		{
			name: "Default values",
			args: []string{"cmd"},
			expectedConfig: Config{
				ShowVersion:         false,
				LogLevel:            "info",
				RateLimitMultiplier: 500,
			},
		},
		{
			name: "Custom log level",
			args: []string{"cmd", "-loglevel", "debug"},
			expectedConfig: Config{
				ShowVersion:         false,
				LogLevel:            "debug",
				RateLimitMultiplier: 500,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			os.Args = tc.args
			config := parseFlags()
			assert.Equal(t, tc.expectedConfig, config)
		})
	}
}

func TestSetupLogging(t *testing.T) {
	testCases := []struct {
		name        string
		logLevel    string
		expectError bool
	}{
		{"Valid log level", "info", false},
		{"Invalid log level", "invalid", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := setupLogging(tc.logLevel)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, logging.INFO, logging.GetLogger().GetLogLevel())
			}
		})
	}
}

func TestSetupDatabase(t *testing.T) {
	testCases := []struct {
		name        string
		dbURL       string
		expectError bool
	}{
		{"Valid DB URL", os.Getenv("DB_URL"), false},
		{"Invalid DB URL", "invalid-url", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbpool, err := setupDatabase(tc.dbURL)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, dbpool)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dbpool)
				dbpool.Close()
			}
		})
	}
}
