package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	config "project/config/auth"
	storage "project/internal/auth/adapter/postgres"
	"project/internal/auth/controller/http"
	"project/internal/auth/service"

	"project/pkg/blacklist"
	"project/pkg/httpserver"
	"project/pkg/postgres"
	"project/pkg/redis"
	"project/pkg/token"

	"github.com/go-chi/chi/v5"
)

func Run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	pgPool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}

	userStorage := storage.New(pgPool.Pool)

	redisClient, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		return fmt.Errorf("error connection to redis: %w", err)
	}

	tokenManager := token.New(cfg.Token)
	blackListManager := blacklist.New(redisClient.Client)

	authService := service.New(service.Params{
		Store:           userStorage,
		Token:           tokenManager,
		Blacklist:       blackListManager,
		TokenExpiration: cfg.Expiration,
	})

	r := chi.NewRouter()
	http.AuthRouter(r, authService, tokenManager, blackListManager)
	httpServer := httpserver.New(r, cfg.HTTP)

	slog.Info("App started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig // wait signal

	slog.Info("App got signal to stop")

	// Controllers close
	httpServer.Close()

	// Adapters close
	redisClient.Close()
	pgPool.Close()

	slog.Info("App stopped!")

	return nil
}
