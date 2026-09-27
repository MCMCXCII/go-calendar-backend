package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	config "project/config/events"
	"project/internal/events/adapter/cache"
	"project/internal/events/adapter/storage"
	"project/internal/events/controller/http"
	"project/internal/events/service"
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

	pgStore := storage.New(pgPool.Pool)

	redisClient, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		return fmt.Errorf("error connection to redis: %w", err)
	}
	blacklistRedisClient, err := redis.New(ctx, redis.Config{URL: cfg.Blacklist.URL})
	if err != nil {
		return fmt.Errorf("error connection to blacklist redis: %w", err)
	}

	cachingStore := cache.New(cache.Params{
		Next:   pgStore,
		Client: redisClient.Client,
		Config: cfg.Cache,
	})

	tokenManager := token.New(cfg.Token)
	blackListManager := blacklist.New(blacklistRedisClient.Client)

	eventService := service.New(service.Params{Store: cachingStore})

	httpMetrics := metrics.NewHTTPServer()

	r := chi.NewRouter()
	http.EventsRouter(r, eventService, tokenManager, blackListManager, httpMetrics)
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

	redisClient.Close()
	pgPool.Close()

	logger.Log.Info("server stopped")

	return nil
}
