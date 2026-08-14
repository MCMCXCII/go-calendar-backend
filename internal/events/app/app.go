package app

import (
	"context"
	"fmt"
	"project/internal/events/cache"
	"project/internal/events/config"
	"project/internal/events/service"
	"project/internal/events/storage"
	"project/internal/events/transport/httpserver"
	"project/internal/platform/blacklist"
	"project/internal/platform/postgres"
	"project/internal/platform/redis"
	"project/internal/platform/token"
)

func Run(ctx context.Context) error {
	cfg := config.MustLoad()

	db, err := postgres.New(ctx, postgres.Params{URL: cfg.Database.URL})
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}
	defer db.Close()

	redisClient, err := redis.New(ctx, redis.Params{URL: cfg.Redis.URL})
	if err != nil {
		return fmt.Errorf("error connection to redis: %w", err)
	}
	defer redisClient.Close()

	pgStore := storage.New(storage.Params{Pool: db.Pool()})

	cachingStore := cache.New(cache.Params{
		Next:  pgStore,
		Redis: redisClient.Client(),
		TTL:   cfg.Cache.TTL,
	})

	tokenParser := token.New(token.Params{Secret: cfg.JWT.SecretKey})
	blacklistChecker := blacklist.New(blacklist.Params{Client: redisClient.Client()})

	eventService := service.New(service.Params{Store: cachingStore})

	server := httpserver.New(httpserver.Params{
		Addr:      cfg.HTTPServer.Address,
		App:       eventService,
		Token:     tokenParser,
		Blacklist: blacklistChecker,
	})

	if err := server.Run(ctx); err != nil {
		return fmt.Errorf("error run HTTP server: %w", err)
	}
	return nil
}
