package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kurochkinivan/user_service/internal/app"
	"github.com/kurochkinivan/user_service/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := app.New(ctx, log, cfg)
	application.Run(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop

	log.Info("stopping application", slog.String("signal", sig.String()))
	application.Stop()
	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		)
	}

	return log
}
