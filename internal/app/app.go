package app

import (
	"cerberus/internal/app/grpc"
	"cerberus/internal/config"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storage config.Storage,
	tokenTTL time.Duration,
) *App {
	// TODO: init storage

	// TODO: init auth service
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
