package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Load[T any]() (*T, error) {
	_ = godotenv.Load()

	var cfg T
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return &cfg, nil
}

type Postgres struct {
	DSN         string `env:"DATABASE_URL,required"`
	AutoMigrate bool   `env:"AUTO_MIGRATE" envDefault:"true"`
}

type Redis struct {
	Addr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

type JWT struct {
	Secret     string        `env:"JWT_SECRET,required"`
	AccessTTL  time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`
	RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" envDefault:"720h"`
	Issuer     string        `env:"JWT_ISSUER" envDefault:"adventist-stack"`
}

type SMTP struct {
	Host string `env:"SMTP_HOST" envDefault:"localhost"`
	Port int    `env:"SMTP_PORT" envDefault:"1025"`
	User string `env:"SMTP_USER"`
	Pass string `env:"SMTP_PASS"`
	From string `env:"SMTP_FROM" envDefault:"no-reply@adventist-stack.local"`
}
