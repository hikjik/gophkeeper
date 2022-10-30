package pg

import (
	"database/sql"
	"embed"
	"errors"

	m "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	// Register some db stuff
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
)

//go:embed migrations/*.sql
var fs embed.FS

type postgresStorage struct {
	db *sql.DB
}

var _ storage.Storage = (*postgresStorage)(nil)

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

// New возвращает объект postgresStorage, реализующий интерфейс Storage
func New(databaseURL string) (storage.Storage, error) {
	if err := migrate(databaseURL); err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	return &postgresStorage{db: db}, nil
}
