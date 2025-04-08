package authrpc

import (
	"authservice/internal/domain/models"
	"context"

	ssov1 "github.com/hoptdev/sso_protos/gen/go/sso"
	"google.golang.org/grpc"
)

type Auth interface {
	Login(ctx context.Context, password string) (models.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (models.TokenPair, error)
	Validate(ctx context.Context, refreshToken string) (bool, int, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (
	*ssov1.LoginResponse, error) {

	t, err := s.auth.Login(ctx, req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &ssov1.LoginResponse{
		RefreshToken: t.RefreshToken,
		AccessToken:  t.AccessToken,
	}, nil
}

func (s *serverAPI) Refresh(ctx context.Context, req *ssov1.RefreshRequest) (
	*ssov1.RefreshResponse, error) {

	t, err := s.auth.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &ssov1.RefreshResponse{
		RefreshToken: t.RefreshToken,
	}, nil
}

func (s *serverAPI) ValidateToken(ctx context.Context, req *ssov1.ValidateTokenRequest) (
	*ssov1.ValidateTokenResponse, error) {

	t, err := s.auth.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &ssov1.RefreshResponse{
		RefreshToken: t.RefreshToken,
	}, nil
}
