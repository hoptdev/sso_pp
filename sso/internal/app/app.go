package app

import (
	grpcapp "authservice/internal/app/grpc"
	"authservice/internal/services/auth"
	"authservice/internal/storage/psql"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, port int, psql_connect string) *App {
	storage, err := psql.New(log, psql_connect)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage)

	grpcserver := grpcapp.New(log, port, authService)

	return &App{
		GRPCServer: grpcserver,
	}
}
