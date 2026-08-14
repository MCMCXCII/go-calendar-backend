package app

import (
	"context"
	"fmt"
	"project/internal/auth/config"
	"project/internal/auth/service"
	"project/internal/auth/storage"
	"project/internal/auth/transport/httpserver"
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

	userStorage := storage.New(storage.Params{Pool: db.Pool()})

	redisClient, err := redis.New(ctx, redis.Params{URL: cfg.Redis.URL})
	if err != nil {
		return fmt.Errorf("error connection to redis: %w", err)
	}
	defer redisClient.Close()

	tokenManager := token.New(token.Params{Secret: cfg.JWT.SecretKey})
	blackListManager := blacklist.New(blacklist.Params{Client: redisClient.Client()})

	authService := service.New(service.Params{
		Store:           userStorage,
		Token:           tokenManager,
		Blacklist:       blackListManager,
		TokenExpiration: cfg.JWTExpiration,
	})

	server := httpserver.New(httpserver.Params{
		Addr:      cfg.HTTPServer.Address,
		App:       authService,
		Token:     tokenManager,
		Blacklist: blackListManager,
	},
	)
	if err := server.Run(ctx); err != nil {
		return fmt.Errorf("error run HTTP server: %w", err)
	}
	return nil
}
