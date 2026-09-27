package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	config "project/config/auth"
	storage "project/internal/auth/adapter/postgres"
	"project/internal/auth/controller/http"
	"project/internal/auth/service"
	"project/pkg/blacklist"
	"project/pkg/httpserver"
	"project/pkg/logger"
	"project/pkg/metrics"
	"project/pkg/postgres"
	"project/pkg/redis"
	"project/pkg/token"
)

func Run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err = logger.New(cfg.Logger); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	pgPool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}
	defer pgPool.Close()

	userStorage := storage.New(pgPool.Pool)

	redis, err := redis.New(ctx, redis.Config{URL: cfg.Blacklist.URL})
	if err != nil {
		return fmt.Errorf("error connection to redis: %w", err)
	}
	defer redis.Close()

	tokenManager := token.New(cfg.Token)
	blackListManager := blacklist.New(redis.Client)

	authService := service.New(service.Params{
		Store:           userStorage,
		Token:           tokenManager,
		Blacklist:       blackListManager,
		TokenExpiration: cfg.Expiration,
	})

	httpMetrics := metrics.NewHTTPServer()

	r := chi.NewRouter()
	http.AuthRouter(r, authService, tokenManager, blackListManager, httpMetrics)
	httpServer := httpserver.New(r, cfg.HTTP)

	logger.Log.Info("server started",
		zap.String("app", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("address", cfg.HTTP.Address),
	)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	logger.Log.Info("shutdown signal received")

	httpServer.Close()

	redis.Close()
	pgPool.Close()

	logger.Log.Info("server stopped")

	return nil
}
