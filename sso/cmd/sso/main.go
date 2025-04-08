package main

import (
	"authservice/internal/app"
	"authservice/internal/config"
	"log/slog"
	"os"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.Load()

	log := setupLogger(cfg.Env)

	log.Info("start application", slog.String("mode", cfg.Env))

	application := app.New(log, cfg.GRPC.Port, cfg.PSQL_Connect)
	application.GRPCServer.Run()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
