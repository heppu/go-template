// Package store abstracts the database layer behind store.Store methods.
package store

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-srvc/mods/sqlxmod"
	"github.com/go-srvc/srvc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Store wraps the sqlxmod module and provides an interface to interact with the database.
type Store struct {
	srvc.Module
	db NamedDB
}

func New(opts ...sqlxmod.Opt) *Store {
	s := &Store{}
	s.Module = sqlxmod.New(append(opts, setDB(s))...)
	return s
}

func (s *Store) Healthy(ctx context.Context) error {
	t := time.Time{}
	if err := s.db.GetContext(ctx, &t, "SELECT NOW()"); err != nil {
		return err
	}
	slog.Info("DB healthy", slog.Time("time_from_db", t))
	return nil
}
