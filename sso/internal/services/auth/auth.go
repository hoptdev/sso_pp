package auth

import (
	"authservice/internal/domain/models"
	"authservice/internal/lib/jwt"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrTokenInvalid   = errors.New("token invalid")
	ErrorUpdateFailed = errors.New("update token failed")

	ErrLoginExists    = errors.New("login exists")
	ErrRegisterFailed = errors.New("register fail")
)

type Auth struct {
	log          *slog.Logger
	userUpdater  UserUpdater
	userProvider UserProvider
	userInserter UserInserter
}

type UserProvider interface {
	GetUserByToken(ctx context.Context, t models.TokenPair) (*models.User, error)
	GetUserByPassword(ctx context.Context, login string, password string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}

type UserUpdater interface {
	UpdateUserToken(ctx context.Context, user *models.User, t models.TokenPair) error
}

type UserInserter interface {
	CreateUser(ctx context.Context, login string, password string) (int, error)
}

func New(log *slog.Logger, userSaver UserUpdater, userProvider UserProvider, userInserter UserInserter) *Auth {
	return &Auth{log, userSaver, userProvider, userInserter}
}

func (a *Auth) Login(ctx context.Context, login string, password string) (models.TokenPair, error) {
	var token models.TokenPair

	user, err := a.userProvider.GetUserByPassword(ctx, login, password)
	if err != nil {
		return token, ErrUserNotFound
	}

	token, err = jwt.NewPair(jwt.UserClaims{UserId: user.Id, Role: 1})
	if err != nil {
		return token, err
	}

	err = a.userUpdater.UpdateUserToken(ctx, user, token)

	if err != nil {
		return token, ErrorUpdateFailed
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, login string, password string) (bool, error) {
	user, err := a.userProvider.GetUserByLogin(ctx, login)
	if err != nil || user != nil {
		return false, ErrLoginExists
	}

	_, err = a.userInserter.CreateUser(ctx, login, password)
	if err != nil {
		return false, ErrRegisterFailed
	}

	return true, nil
}

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (models.TokenPair, error) {
	var token models.TokenPair
	token.RefreshToken = refreshToken

	user, _ := a.userProvider.GetUserByToken(ctx, token)
	if user == nil {
		return token, ErrUserNotFound
	}

	res, claims, err := jwt.ValidateToken(token)

	if !res && err != nil {
		return token, ErrTokenInvalid
	}

	ntoken, err := jwt.NewPair(claims)
	if err != nil {
		return token, ErrTokenInvalid
	}

	err = a.userUpdater.UpdateUserToken(ctx, user, ntoken)
	if err != nil {
		return token, ErrorUpdateFailed
	}

	return ntoken, nil
}

func (a *Auth) Validate(ctx context.Context, refreshToken string) (isValid bool, userId int, err error) {
	a.log.Info(fmt.Sprintf("Validate %v", refreshToken))

	var token models.TokenPair
	token.RefreshToken = refreshToken

	user, _ := a.userProvider.GetUserByToken(ctx, token)
	if user == nil {
		return false, userId, ErrUserNotFound
	}

	res, _, err := jwt.ValidateToken(token)

	if !res && err != nil {
		return false, userId, ErrTokenInvalid
	}

	return true, user.Id, nil
}
