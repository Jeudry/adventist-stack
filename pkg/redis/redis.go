// Package redis provee un cliente Redis compartido, usado tanto para cache
// como para pub/sub de notificaciones en tiempo real.
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client envuelve *redis.Client para exponer helpers de alto nivel.
type Client struct {
	*redis.Client
}

// Connect crea el cliente y verifica conectividad.
func Connect(ctx context.Context, addr, password string, db int) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := rdb.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("redis: ping: %w", err)
	}
	return &Client{rdb}, nil
}
