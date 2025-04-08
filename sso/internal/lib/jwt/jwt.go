package jwt

import (
	"authservice/internal/domain/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserId int
	Role   int
	jwt.RegisteredClaims
}

const test_key = "L0sUZ3a1m8XzuGRus3l9wnhIMwSc6cDzFJNnWwFZRMY="

const (
	refreshMultiple = 604800
	accessMultiple  = 300
)

/*func GenerateAndValidateToken() (models.TokenPair, error) {
	result, claims, nerr := ValidateToken(token)
	if !result || nerr != nil {
		return token, errors.New("invalid refresh token")
	}

	token.NewTokens(claims)

	return token, nil
} */

func NewPair(claims UserClaims) (models.TokenPair, error) {
	var token models.TokenPair

	refresh, err := NewRefreshToken(claims)
	access, err2 := NewAccessToken(claims)
	if err != nil || err2 != nil {
		return token, errors.New("generate tokens failed")
	}

	token.RefreshToken = refresh
	token.AccessToken = access

	return token, nil
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

func ValidateToken(t models.TokenPair) (bool, UserClaims, error) {
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
