package suite

import (
	"cerberus/internal/config"
	"context"
	cerberusv1 "github.com/alexbro4u/contracts/gen/go/cerberus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient cerberusv1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	var cfg config.Config
	config.LoadConfig("../config/local_tests.yaml", &cfg)

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.ReadTimeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(
		context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial gRPC: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        &cfg,
		AuthClient: cerberusv1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg config.Config) string {
	return net.JoinHostPort(cfg.GRPC.Host, strconv.Itoa(cfg.GRPC.Port))
}
