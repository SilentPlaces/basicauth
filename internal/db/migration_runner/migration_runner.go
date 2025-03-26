package migration_runner

import (
	"database/sql"
	"github.com/SilentPlaces/basicauth.git/pkg/constants"
	"github.com/google/wire"
	"github.com/pressly/goose/v3"
)

// To run migration from codes. but it will be unused. i just wrote the codes
type MigrationRunner interface {
	migrationUp() error
	migrationDown() error
}

type migrationRunner struct {
	db *sql.DB
}

func NewMigrationRunner(db *sql.DB) MigrationRunner {
	return &migrationRunner{db: db}
}

func (m *migrationRunner) migrationUp() error {
	return goose.Up(m.db, constants.MigrationsDirectory)
}

func (m *migrationRunner) migrationDown() error {
	return goose.Down(m.db, constants.MigrationsDirectory)
}

var ProviderSet = wire.NewSet(NewMigrationRunner)
