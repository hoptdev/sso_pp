package psql

import (
	"authservice/internal/domain/models"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
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

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	query := "SELECT id, token FROM users WHERE login = $1 LIMIT 1"
	row := conn.QueryRow(ctx, query, login)
	err = row.Scan(&user.Id, &user.Token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *Storage) GetUserByPassword(ctx context.Context, login string, password string) (*models.User, error) {
	var user models.User
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	query := "SELECT id, token FROM users WHERE pass = $1 AND login = $2"
	row := conn.QueryRow(ctx, query, password, login)
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

func (s *Storage) CreateUser(ctx context.Context, login string, password string) (int, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return 0, err
	}

	defer conn.Release()

	query := "INSERT INTO users (login, pass, token) VALUES ($1, $2, $3) RETURNING id;"
	row := conn.QueryRow(ctx, query, login, password, "")
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
