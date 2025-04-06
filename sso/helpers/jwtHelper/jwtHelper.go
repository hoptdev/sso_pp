package jwtHelper

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const test_key = "L0sUZ3a1m8XzuGRus3l9wnhIMwSc6cDzFJNnWwFZRMY="

const (
	refreshMultiple = 604800
	accessMultiple  = 300
)

type UserClaims struct {
	UserId int
	Role   int
	jwt.RegisteredClaims
}

type TokenPair struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
}

func (t TokenPair) ToJson() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (t TokenPair) String() string {
	return t.RefreshToken + " " + t.AccessToken
}

func (token TokenPair) GenerateAndValidateToken() (TokenPair, error) {
	result, claims, nerr := ValidateToken(token)
	if !result || nerr != nil {
		return token, errors.New("invalid refresh token")
	}

	token.NewTokens(claims)

	return token, nil
}

func (token *TokenPair) NewTokens(claims UserClaims) error {
	refresh, err := NewRefreshToken(claims)
	access, err2 := NewAccessToken(claims)
	if err != nil || err2 != nil {
		return errors.New("generate tokens failed")
	}

	token.RefreshToken = refresh
	token.AccessToken = access

	return nil
}

func NewRefreshToken(u UserClaims) (string, error) {
	return newToken(u, refreshMultiple)
}

func NewAccessToken(u UserClaims) (string, error) {
	return newToken(u, accessMultiple)
}

func newToken(u UserClaims, seconds int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userId": u.UserId, "role": u.Role, "exp": jwt.NewNumericDate(time.Now().Add(time.Duration(seconds) * 1e9))})

	tokenStr, err := token.SignedString([]byte(test_key))

	return tokenStr, err
}

func ValidateToken(t TokenPair) (bool, UserClaims, error) {
	var claims UserClaims
	parser, err := jwt.ParseWithClaims(t.RefreshToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(test_key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		//log.Print(err)
		return false, claims, err
	}

	return parser.Valid, claims, err
}
