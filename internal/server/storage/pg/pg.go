package pg

import (
	"embed"
	"errors"

	m "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	// Register some db stuff
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//go:embed migrations/*.sql
var fs embed.FS

func migrate(databaseURL string) error {
	sourceDriver, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	migrateInstance, err := m.NewWithSourceInstance("iofs", sourceDriver, databaseURL)
	if err != nil {
		return err
	}
	err = migrateInstance.Up()
	if err != nil && !errors.Is(err, m.ErrNoChange) {
		return err
	}
	return nil
}
