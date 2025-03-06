package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/elisasre/go-common/v2/sentryutil"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

func RunMigrations(db *sql.DB) error {
	source, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration fs from embedded fs: %w", err)
	}

	span := sentryutil.MakeSpan(context.Background(), 1)
	defer span.Finish()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return err
	}

	dbInfo(m.Version())
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("already at latest migration")
		return nil
	}

	dbInfo(m.Version())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func dbInfo(version uint, dirty bool, err error) {
	l := slog.With(slog.Uint64("current_version", uint64(version)))
	l = l.With(slog.Bool("dirty", dirty))
	if err != nil {
		l = l.With(slog.String("error", err.Error()))
	}
	l.Info("DB info")
}
