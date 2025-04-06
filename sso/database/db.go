package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DbPool *pgxpool.Pool

const connect string = "postgres://postgres:admin@localhost:5432/authService_pp4"

func SetupDatabase() {
	log.Println("Connect to db")
	config, err := pgxpool.ParseConfig(connect)
	if err != nil {
		log.Fatal(err)
	}

	config.MaxConns = 50
	config.MinConns = 10
	config.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DbPool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
}
