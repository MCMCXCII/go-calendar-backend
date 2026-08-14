package config

import (
	"project/internal/platform/config"
	"time"
)

type Config struct {
	Env           string `env:"ENV" env-default:"local"`
	HTTPServer    config.HTTPServer
	JWT           config.JWT
	JWTExpiration time.Duration `env:"JWT_EXPIRATION" env-default:"24h"`
	Redis         config.Redis
	Database      Database
}

type Database struct {
	URL string `env:"AUTH_DB_URL" env-required:"true"`
}

func MustLoad() *Config {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	return cfg
}
