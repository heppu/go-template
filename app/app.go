package app

//go:generate go tool github.com/ogen-go/ogen/cmd/ogen --clean --config ../.ogen.yaml --target ../api -package api ../openapi.yaml

import (
	"context"
	"log/slog"

	"github.com/heppu/go-template/api"
)

type Store interface {
	Healthy(context.Context) error
}

type App struct {
	s Store
}

func New(s Store) *App {
	return &App{s: s}
}

func (a *App) Healthz(ctx context.Context) (*api.Healthy, error) {
	if err := a.s.Healthy(ctx); err != nil {
		return nil, err
	}
	return &api.Healthy{Message: "OK"}, nil
}

// NewError can be used to provide custom error responses based on the error.
func (a *App) NewError(ctx context.Context, err error) *api.ErrorRespStatusCode {
	slog.Error("internal server error", slog.String("err", err.Error()))
	return &api.ErrorRespStatusCode{
		StatusCode: 500,
		Response: api.Error{
			Error: "internal server error",
		},
	}
}
