package store

import (
	"context"

	"github.com/go-srvc/mods/sqlxmod"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	*sqlxmod.DB
}

func New(opts ...sqlxmod.Opt) *Store {
	return &Store{DB: sqlxmod.New(opts...)}
}

func (s *Store) Healthy(ctx context.Context) error {
	return s.DB.DB().PingContext(ctx)
}
