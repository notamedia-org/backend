package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	Port           int    `env:"PORT"    envDefault:"8887"`
	PgConnection   string `env:"PG_CONNECTION_STRING" envDefault:"http://postgres@root:root:3005/icpc"`
	MigrationsPath string `env:"MIGRATIONS_PATH" envDefault:"./migrations"`
	JwtSecret      string `env:"JWT_SECRET" envDefault:"root"`
}

func ReadFromEnv() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("cannot load env from environment: %w", err)
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}

	return &cfg, nil
}
