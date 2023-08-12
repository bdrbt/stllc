package database

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

//go:embed migrations
var migrations embed.FS

func MigrateDB(dbURL string) error {
	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("cannot create iofs:%w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		return fmt.Errorf("cannot create migration instance:%w", err)
	}

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("cannot migrate database:%w", err)
	}

	return nil
}
