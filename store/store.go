package store

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-srvc/mods/sqlxmod"
	"github.com/go-srvc/srvc"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	srvc.Module
	db *sqlx.DB
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
