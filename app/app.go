// Package app wraps the storage layer with application logic and provides an interface for API layer.
package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-srvc/srvc"
	"github.com/heppu/go-template/api"
)

const ErrServiceNotHealthy = srvc.ErrStr("service not healthy")

type Store interface {
	Healthy(context.Context) error
}

type App struct {
	api.UnimplementedHandler
	s Store
}

func New(s Store) *App {
	return &App{s: s}
}

func (a *App) Healthz(ctx context.Context) (*api.Healthy, error) {
	if err := a.s.Healthy(ctx); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceNotHealthy, err)
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
