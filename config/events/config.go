package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"project/internal/events/adapter/cache"
	"project/pkg/httpserver"
	"project/pkg/postgres"
	"project/pkg/redis"
	"project/pkg/token"
)

type App struct {
	Name    string `envconfig:"APP_NAME"    default:"events-service"`
	Version string `envconfig:"APP_VERSION" default:"dev"`
}

type Config struct {
	App      App
	HTTP     httpserver.Config
	Token    token.Config
	Redis    redis.Config
	Postgres postgres.Config
	Cache    cache.Config
}

func New() (Config, error) {
	var config Config

	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		return config, fmt.Errorf("godotenv.Load: %w", err)
	}

	if err := envconfig.Process("EVENTS", &config); err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
