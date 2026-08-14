package config

import "time"

type HTTPServer struct {
	Address     string        `env:"HTTP_ADDRESS" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `env:"HTTP_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type Redis struct {
	URL string `env:"REDIS_URL" env-required:"true"`
}

type JWT struct {
	SecretKey string `env:"JWT_SECRET" env-required:"true"`
}
