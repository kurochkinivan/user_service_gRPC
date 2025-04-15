package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kurochkinivan/user_service/internal/config"
	usergrpc "github.com/kurochkinivan/user_service/internal/controller/grpc/users"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
	timeout    time.Duration
}

func New(log *slog.Logger, cfg config.GRPCConfig, user usergrpc.User) *App {
	server := grpc.NewServer(grpc.ConnectionTimeout(cfg.Timeout))

	validate := validator.New(validator.WithRequiredStructEnabled())

	usergrpc.Register(server, validate, user)

	return &App{
		log:        log,
		gRPCServer: server,
		port:       cfg.Port,
		timeout:    cfg.Timeout,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("op", op)).
		Info("stopping grpc server...", slog.String("port", a.port))

	a.gRPCServer.GracefulStop()
}
