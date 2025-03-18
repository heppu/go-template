package store_test

import (
	"context"
	"testing"

	"github.com/go-srvc/mods/sqlxmod"
	"github.com/heppu/go-template/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamedDB_NamedGetContext(t *testing.T) {
	s := sqlxmod.New(store.WithOtel())
	require.NoError(t, s.Init())
	db := store.NamedDB{s.DB()}

	const expected = "bar"
	input := struct{ Foo string }{Foo: expected}
	var got string

	err := db.NamedGetContext(context.Background(), &got, "SELECT :foo", input)
	assert.NoError(t, err)

	assert.Equal(t, expected, got)
	require.NoError(t, s.Stop())
}
