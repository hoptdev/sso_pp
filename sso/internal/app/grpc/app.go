package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	authrpc "sso/internal/grpc/auth"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()
	authrpc.Register(gRPCServer)

	return &App{log, gRPCServer, port}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"

	log := app.log.With(slog.String("op", op), slog.Int("port", app.port))

	log.Info("starting gRPC server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gprc server is running", slog.String("address", l.Addr().String()))

	if err := app.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) Stop() {
	const op = "grpcapp.Stop"

	app.log.With(slog.String("op", op)).Info("grpc server stop")

	app.gRPCServer.GracefulStop()
}
