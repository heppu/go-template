package store

import (
	"fmt"
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type Opt func(*Store) error

func FromEnv() Opt {
	return func(s *Store) error {
		port := os.Getenv("POSTGRES_PORT")
		user := os.Getenv("POSTGRES_USER")
		host := os.Getenv("POSTGRES_HOST")
		dbName := os.Getenv("POSTGRES_DB")
		passwd := os.Getenv("POSTGRES_PASSWORD")
		sslMode := os.Getenv("POSTGRES_SSLMODE")

		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, passwd, host, port, dbName, sslMode)
		db, err := otelsqlx.Open("pgx", dsn, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
		if err != nil {
			return fmt.Errorf("failed to open db using env variables: %w", err)
		}

		s.db = db
		return nil
	}
}
