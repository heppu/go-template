package store

import (
	"context"
	"fmt"

	"github.com/go-tstr/tstr/strerr"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const ErrDBNotSet = strerr.Error("db not set")

type Store struct {
	opts []Opt
	db   *sqlx.DB
	done chan struct{}
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

func (s *Store) Healthy(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
