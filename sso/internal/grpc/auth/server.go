package authrpc

import (
	"context"

	ssov1 "github.com/hoptdev/sso_pp/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (
	*ssov1.LoginResponse, error) {
	panic("")
}
