package store

import (
	"cmp"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/XSAM/otelsql"
	"github.com/go-srvc/mods/sqlxmod"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

//go:embed migrations/*.sql
var fs embed.FS

func WithDefaults() []sqlxmod.Opt {
	return []sqlxmod.Opt{
		WithOtel(),
		WithMigrations(),
	}
}

func WithOtel() sqlxmod.Opt {
	return func(db *sqlxmod.DB) error {
		port := get("POSTGRES_PORT", "5432")
		user := get("POSTGRES_USER", "root")
		host := get("POSTGRES_HOST", "127.0.0.1")
		dbName := get("POSTGRES_DB", "root")
		passwd := get("POSTGRES_PASSWORD", "root")
		sslMode := get("POSTGRES_SSLMODE", "disable")
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, passwd, host, port, dbName, sslMode)
		return sqlxmod.WithOtel("pgx", dsn, otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL))(db)
	}
}

func WithMigrations() sqlxmod.Opt {
	return func(db *sqlxmod.DB) error {
		source, err := iofs.New(fs, "migrations")
		if err != nil {
			return fmt.Errorf("failed to create migration fs from embedded fs: %w", err)
		}
		driver, err := postgres.WithInstance(db.DB().DB, &postgres.Config{})
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
}

func setDB(s *Store) sqlxmod.Opt {
	return func(db *sqlxmod.DB) error {
		s.db = NamedDB{db.DB()}
		return nil
	}
}

func dbInfo(version uint, dirty bool, err error) {
	l := slog.With(slog.Uint64("current_version", uint64(version)))
	l = l.With(slog.Bool("dirty", dirty))
	if err != nil {
		l = l.With(slog.String("error", err.Error()))
	}
	l.Info("DB info")
}

func get(key, defVal string) string {
	return cmp.Or(os.Getenv(key), defVal)
}
