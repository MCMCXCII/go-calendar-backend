package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Config struct {
	Mode  string `envconfig:"LOG_MODE" default:"dev"`
	Level string `envconfig:"LOG_LEVEL" default:"info"`
}

var Log *zap.Logger = zap.NewNop()

func New(c Config) error {
	level, err := zap.ParseAtomicLevel(c.Level)
	if err != nil {
		return fmt.Errorf("parse log level %q: %w", c.Level, err)
	}

	var cfg zap.Config
	if c.Mode == "json" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.Level = level

	zl, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("build zap logger: %w", err)
	}

	Log = zl
	return nil
}
