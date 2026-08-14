package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func Load[T any]() (*T, error) {
	var cfg T

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return &cfg, nil
}
