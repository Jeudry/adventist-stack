// Package config carga la configuración desde variables de entorno hacia
// structs tipadas y valida al arranque, de modo que la app falle rápido si
// falta una credencial en lugar de romperse en runtime.
//
// Ningún secreto se hardcodea: todo viene del entorno (.env en desarrollo,
// variables reales inyectadas en producción).
package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Load lee el archivo .env (si existe) y parsea el entorno en T.
// En producción el .env no existe y se usan las variables inyectadas.
func Load[T any]() (*T, error) {
	// El error de godotenv se ignora a propósito: en producción no hay .env.
	_ = godotenv.Load()

	var cfg T
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return &cfg, nil
}

// Postgres agrupa la conexión a la base de datos.
type Postgres struct {
	DSN         string `env:"DATABASE_URL,required"`
	AutoMigrate bool   `env:"AUTO_MIGRATE" envDefault:"true"`
}

// Redis agrupa la conexión a Redis (cache + pub/sub de notificaciones).
type Redis struct {
	Addr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

// JWT agrupa los secretos y tiempos de vida de los tokens.
type JWT struct {
	Secret     string        `env:"JWT_SECRET,required"`
	AccessTTL  time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`
	RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" envDefault:"720h"`
	Issuer     string        `env:"JWT_ISSUER" envDefault:"adventist-stack"`
}

// SMTP agrupa la configuración del servidor de correo.
type SMTP struct {
	Host string `env:"SMTP_HOST" envDefault:"localhost"`
	Port int    `env:"SMTP_PORT" envDefault:"1025"`
	User string `env:"SMTP_USER"`
	Pass string `env:"SMTP_PASS"`
	From string `env:"SMTP_FROM" envDefault:"no-reply@adventist-stack.local"`
}
