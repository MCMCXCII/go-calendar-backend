package config

import (
	"project/internal/platform/config"
	"time"
)

type Config struct {
	Env        string `env:"ENV" env-default:"local"`
	HTTPServer config.HTTPServer
	Database   Database
	Redis      config.Redis
	JWT        config.JWT
	Cache      Cache
}

type Cache struct {
	TTL time.Duration `env:"CACHE_TTL" env-default:"300s"`
}

type Database struct {
	URL string `env:"EVENTS_DB_URL" env-required:"true"`
}

func MustLoad() *Config {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	return cfg
}
