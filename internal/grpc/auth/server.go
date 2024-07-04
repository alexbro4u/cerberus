package auth

import (
	"context"
	cerberusv1 "github.com/alexbro4u/contracts/gen/go/cerberus"
	"google.golang.org/grpc"
)

type serverApi struct {
	cerberusv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	cerberusv1.RegisterAuthServer(gRPC, &serverApi{})
}

func (s *serverApi) Login(
	ctx context.Context,
	req *cerberusv1.LoginRequest,
) (*cerberusv1.LoginResponse, error) {

	return &cerberusv1.LoginResponse{
		Token: "token",
	}, nil
}

func (s *serverApi) Register(
	ctx context.Context,
	req *cerberusv1.RegisterRequest,
) (*cerberusv1.RegisterResponse, error) {
	panic("not implemented")
}

func (s *serverApi) IsAdmin(
	ctx context.Context,
	req *cerberusv1.IsAdminRequest,
) (*cerberusv1.IsAdminResponse, error) {
	panic("not implemented")
}
