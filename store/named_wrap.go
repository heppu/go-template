package store

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// NamedDB is a wrapper around sqlx.DB that adds NamedGetContext method.
// https://github.com/jmoiron/sqlx/issues/154#issuecomment-148216948
type NamedDB struct {
	*sqlx.DB
}

// NamedGetContext executes a named query with the given context and arguments.
// It prepares the statement, executes it, and scans the result into the destination.
// This is commonly used for INSERT ... RETURNING query.
func (n NamedDB) NamedGetContext(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	stmt, err := n.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, dest, arg)
}
