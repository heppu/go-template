package store

import (
	"fmt"
	"os"

	"github.com/go-tstr/tstr/strerr"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const ErrDBNotSet = strerr.Error("db not set")

type Opt func(*Store) error

type Store struct {
	opts []Opt
	db   *sqlx.DB
	done chan struct{}
}

func FromEnv() Opt {
	return func(s *Store) error {
		port := os.Getenv("POSTGRES_PORT")
		user := os.Getenv("POSTGRES_USER")
		host := os.Getenv("POSTGRES_HOST")
		dbName := os.Getenv("POSTGRES_DB")
		passwd := os.Getenv("POSTGRES_PASSWORD")
		sslMode := os.Getenv("POSTGRES_SSLMODE")

		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, passwd, host, port, dbName, sslMode)
		db, err := sqlx.Open("pgx", dsn)
		if err != nil {
			return fmt.Errorf("failed to open db using env variables: %w", err)
		}

		s.db = db
		return nil
	}
}

func New(opts ...Opt) *Store {
	return &Store{
		opts: opts,
		done: make(chan struct{}),
	}
}

func (s *Store) Name() string {
	return "store.Store"
}

func (s *Store) Init() error {
	for _, opt := range s.opts {
		if err := opt(s); err != nil {
			return err
		}
	}
	if s.db == nil {
		return ErrDBNotSet
	}
	if err := RunMigrations(s.db.DB); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func (s *Store) Run() error {
	<-s.done
	return nil
}

func (s *Store) Stop() error {
	close(s.done)
	return s.db.Close()
}

func (s *Store) Healthy() error {
	return s.db.Ping()
}
