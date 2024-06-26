package auth

import (
	cerberusv1 "github.com/alexbro4u/contarcts/gen/go/cerberus"
	"google.golang.org/grpc"
)

type serverApi struct {
	cerberusv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	cerberusv1.RegisterAuthServer(gRPC, &serverApi{})
}
