package grpcapp

import (
	authgrpc "cerberus/internal/grpc/auth"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates a new instance of the gRPC server
func New(
	log *slog.Logger,
	authService authgrpc.Auth,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	reflection.Register(gRPCServer)

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}

}

// MustRun runs the gRPC server and panics if an error occurs
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs the gRPC server
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping grpc server")

	a.gRPCServer.GracefulStop()
}
