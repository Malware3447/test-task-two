package app

import (
	"context"
	"test-task-two/internal/api"
)

type App struct {
	router *api.Router
}

func NewApp(router *api.Router) *App {
	return &App{router: router}
}

func (a *App) Init(ctx context.Context) {
	const op = "app.Init"
	ctx = context.WithValue(ctx, "app", op)

	a.router.Init(ctx)
}
