package psql

import (
	"authservice/internal/domain/models"
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	log    *slog.Logger
	dbPool *pgxpool.Pool
}

func New(log *slog.Logger, connect string) (*Storage, error) {
	config, err := pgxpool.ParseConfig(connect)

	if err != nil {
		return nil, err
	}

	config.MaxConns = 50
	config.MinConns = 10
	config.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Storage{
		log:    log,
		dbPool: dbPool,
	}, nil
}

func (s *Storage) GetUserByToken(ctx context.Context, t models.TokenPair) (*models.User, error) {
	var user models.User
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	query := "SELECT id, token FROM users WHERE token = $1"
	row := conn.QueryRow(ctx, query, t.RefreshToken)
	err = row.Scan(&user.Id, &user.Token)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) GetUserByPassword(ctx context.Context, password string) (*models.User, error) {
	var user models.User
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	query := "SELECT id, token FROM users WHERE pass = $1"
	row := conn.QueryRow(ctx, query, password)
	err = row.Scan(&user.Id, &user.Token)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) UpdateUserToken(ctx context.Context, user *models.User, t models.TokenPair) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	query := "UPDATE users SET token=$1 WHERE id=$2"
	conn.QueryRow(ctx, query, t.RefreshToken, user.Id)

	return nil
}
