package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	config "project/config/events"
	"project/internal/events/adapter/cache"
	"project/internal/events/adapter/storage"
	"project/internal/events/controller/http"
	"project/internal/events/service"
	"project/pkg/blacklist"
	"project/pkg/httpserver"
	"project/pkg/postgres"
	"project/pkg/redis"
	"project/pkg/token"
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

	pgStore := storage.New(pgPool.Pool)

	redisClient, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		return fmt.Errorf("error connection to redis: %w", err)
	}

	cachingStore := cache.New(cache.Params{
		Next:   pgStore,
		Client: redisClient.Client,
		Config: cfg.Cache,
	})

	tokenManager := token.New(cfg.Token)
	blackListManager := blacklist.New(redisClient.Client)

	eventService := service.New(service.Params{Store: cachingStore})

	r := chi.NewRouter()
	http.EventsRouter(r, eventService, tokenManager, blackListManager)
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
