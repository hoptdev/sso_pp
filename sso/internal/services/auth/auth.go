package auth

import (
	"authservice/internal/domain/models"
	"authservice/internal/lib/jwt"
	"context"
	"errors"
	"log/slog"
)

type Auth struct {
	log          *slog.Logger
	userUpdater  UserUpdater
	userProvider UserProvider
}

type UserProvider interface {
	GetUserByToken(ctx context.Context, t models.TokenPair) (*models.User, error)
	GetUserByPassword(ctx context.Context, password string) (*models.User, error)
}

type UserUpdater interface {
	UpdateUserToken(ctx context.Context, user *models.User, t models.TokenPair) error
}

func New(log *slog.Logger, userSaver UserUpdater, userProvider UserProvider) *Auth {
	return &Auth{log, userSaver, userProvider}
}

func (a *Auth) Login(ctx context.Context, password string) (models.TokenPair, error) {
	var token models.TokenPair

	user, err := a.userProvider.GetUserByPassword(ctx, password)
	if err != nil {
		return token, errors.New("user not found")
	}

	token, err = jwt.NewPair(jwt.UserClaims{UserId: user.Id, Role: 1})
	if err != nil {
		return token, err
	}

	err = a.userUpdater.UpdateUserToken(ctx, user, token)

	if err != nil {
		return token, errors.New("update token failed")
	}

	return token, nil
}

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (models.TokenPair, error) {
	var token models.TokenPair
	token.RefreshToken = refreshToken

	user, _ := a.userProvider.GetUserByToken(ctx, token)
	if user == nil {
		return token, errors.New("user not found")
	}

	res, claims, err := jwt.ValidateToken(token)

	if !res && err != nil {
		return token, errors.New("token invalid")
	}

	ntoken, err := jwt.NewPair(claims)
	if err != nil {
		return token, errors.New("fail")
	}

	err = a.userUpdater.UpdateUserToken(ctx, user, ntoken)
	if err != nil {
		return token, errors.New("update token failed")
	}

	return ntoken, nil
}
