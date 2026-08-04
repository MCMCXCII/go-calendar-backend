package app

import (
	"context"
	"fmt"
	"project/internal/auth/config"
	"project/internal/auth/service"
	"project/internal/auth/storage"
	"project/internal/auth/transport/httpserver"
)

func Run(ctx context.Context) error {
	cfg := config.MustLoad()
	userStorage, err := storage.New(ctx, storage.Params{URL: cfg.Database.URL})
	if err != nil {
		return fmt.Errorf("error connection to database: %w", err)
	}

	authService := service.New(service.Params{Store: userStorage, JWTSecret: cfg.Auth.JWT.SecretKey})
	server := httpserver.New(httpserver.Params{Addr: cfg.Auth.HTTPServer.Address, App: authService})
	if err := server.Run(ctx); err != nil {
		return fmt.Errorf("error run HTTP server: %w", err)
	}
	return nil
}
