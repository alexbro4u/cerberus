package app

import (
	"cerberus/internal/app/grpc"
	"cerberus/internal/config"
	"cerberus/internal/services/auth"
	"cerberus/internal/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storageCfg config.Storage,
	tokenTTL time.Duration,
) *App {

	storage, err := postgres.NewStorage(storageCfg.User, storageCfg.Password, storageCfg.DbName, storageCfg.Host)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
