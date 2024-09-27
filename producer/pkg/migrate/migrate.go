package migrate

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/mchatzis/go/producer/pkg/logging"
)

var logger = logging.GetLogger()

type Migrator struct {
	m *migrate.Migrate
}

func NewMigrator() (*Migrator, error) {
	m, err := migrate.New("file:///app/internal/db/migrations", os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}
	return &Migrator{m: m}, nil
}

func (mgr *Migrator) Up() error {
	logger.Info("Migrating up...")
	err := mgr.m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (mgr *Migrator) Down() error {
	logger.Info("Migrating down...")
	err := mgr.m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (mgr *Migrator) Close() error {
	sourceErr, dbErr := mgr.m.Close()
	if sourceErr != nil {
		return fmt.Errorf("source error: %v", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("database error: %v", dbErr)
	}
	return nil
}
