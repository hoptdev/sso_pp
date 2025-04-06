package userHelper

import (
	"context"
	db "main/database"
	"main/helpers/jwtHelper"
)

type User struct {
	Id    int
	Pass  string
	Token string `json:"-"`
}

func GetUserByToken(t jwtHelper.TokenPair, ctx context.Context) (*User, error) {
	var user User
	conn, err := db.DbPool.Acquire(ctx)
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

func GetUserByPassword(ctx context.Context, password string) (*User, error) {
	var user User
	conn, err := db.DbPool.Acquire(ctx)
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

func (u *User) UpdateToken(ctx context.Context, t jwtHelper.TokenPair) error {
	conn, err := db.DbPool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	query := "Update users SET token=$1 WHERE id=$2"
	conn.QueryRow(ctx, query, t.RefreshToken, u.Id)

	return nil
}
