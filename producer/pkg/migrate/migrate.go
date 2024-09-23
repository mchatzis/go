package migrate

import (
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/mchatzis/go/producer/pkg/logging"
)

var logger = logging.GetLogger()

func Migrate() {
	m, err := migrate.New("file:///app/internal/db/migrations", os.Getenv("DB_URL"))
	if err != nil {
		logger.Fatalf("Failed to open connection to migrate with error: %v", err)
	}

	logger.Info("Migrating up...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatalf("Failed to migrate up with error: %v", err)
	}

	logger.Info("Migrating down...")
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		logger.Fatalf("Failed to migrate down with error: %v", err)
	}
}
