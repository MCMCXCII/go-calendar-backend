package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"project/pkg/httpserver"
	"project/pkg/postgres"
	"project/pkg/redis"
	"project/pkg/token"
)

type App struct {
	Name    string `envconfig:"APP_NAME"    default:"auth-service"`
	Version string `envconfig:"APP_VERSION" default:"dev"`
}

type Config struct {
	App        App
	HTTP       httpserver.Config
	Token      token.Config
	Redis      redis.Config
	Postgres   postgres.Config
	Expiration time.Duration `envconfig:"JWT_EXPIRATION" default:"24h"`
}

func New() (Config, error) {
	var config Config

	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		return config, fmt.Errorf("godotenv.Load: %w", err)
	}

	if err := envconfig.Process("AUTH", &config); err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
