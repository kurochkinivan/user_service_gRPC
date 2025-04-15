package app

import (
	"context"
	"log/slog"
	"time"

	grpcapp "github.com/kurochkinivan/user_service/internal/app/grpc"
	pgapp "github.com/kurochkinivan/user_service/internal/app/pg"
	"github.com/kurochkinivan/user_service/internal/config"
	"github.com/kurochkinivan/user_service/internal/usecase/storage/pg"
	"github.com/kurochkinivan/user_service/internal/usecase/user"
)

type App struct {
	log           *slog.Logger
	GRPCApp       *grpcapp.App
	PostgreSQLApp *pgapp.App
}

func New(ctx context.Context, log *slog.Logger, cfg *config.Config) *App {
	pgApp := pgapp.New(ctx, log, cfg.PostgreSQL)

	repository := pg.New(pgApp.Pool)

	userService := user.New(log, repository)

	gRPCApp := grpcapp.New(log, cfg.GRPC, userService)

	return &App{
		GRPCApp:       gRPCApp,
		PostgreSQLApp: pgApp,
		log:           log,
	}
}

func (a *App) Run(ctx context.Context) {
	go a.PostgreSQLApp.MustRun(ctx, 5, 5*time.Second)
	go a.GRPCApp.MustRun()
}

func (a *App) Stop() {
	a.GRPCApp.Stop()
	a.PostgreSQLApp.Stop()
}
