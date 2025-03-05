package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/heppu/go-template/api"
	"github.com/heppu/go-template/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	a := app.New(nil)

	testErr := errors.New("some private test error that should be only logged but not be exposed via API")
	err := a.NewError(context.Background(), testErr)
	require.NotNil(t, err)

	expected := api.ErrorRespStatusCode{
		StatusCode: 500,
		Response: api.Error{
			Error: "internal server error",
		},
	}

	assert.Equal(t, expected, *err)
}
