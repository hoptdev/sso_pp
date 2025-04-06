package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, port int, psql_connect string) *App {

	grpcserver := grpcapp.New(log, port)

	return &App{
		GRPCServer: grpcserver,
	}
}
