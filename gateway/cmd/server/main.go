package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Jeudry/adventist-stack/gateway/internal/clients"
	"github.com/Jeudry/adventist-stack/gateway/internal/handlers"
	"github.com/Jeudry/adventist-stack/gateway/internal/router"
	"github.com/Jeudry/adventist-stack/pkg/config"
	"github.com/Jeudry/adventist-stack/pkg/jwt"
	"github.com/Jeudry/adventist-stack/pkg/logger"
)

type Config struct {
	Env               string        `env:"ENV" envDefault:"dev"`
	HTTPPort          string        `env:"GATEWAY_HTTP_PORT" envDefault:"8080"`
	AuthAddr          string        `env:"AUTH_GRPC_ADDR" envDefault:"localhost:50051"`
	NotificationsAddr string        `env:"NOTIFICATIONS_GRPC_ADDR" envDefault:"localhost:50052"`
	AllowedOrigins    string        `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	RateLimit         int           `env:"RATE_LIMIT_REQUESTS" envDefault:"100"`
	RateWindow        time.Duration `env:"RATE_LIMIT_WINDOW" envDefault:"1m"`
	JWT               config.JWT
}

func main() {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	log := logger.New("gateway", cfg.Env)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	grpcClients, err := clients.New(cfg.AuthAddr, cfg.NotificationsAddr)
	if err != nil {
		log.Error("failed to create gRPC clients", "err", err)
		os.Exit(1)
	}
	defer grpcClients.Close()

	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	handler := router.New(router.Deps{
		JWT:            jwtManager,
		AuthHandler:    handlers.NewAuthHandler(grpcClients.Auth),
		AllowedOrigins: strings.Split(cfg.AllowedOrigins, ","),
		RateLimit:      cfg.RateLimit,
		RateWindow:     cfg.RateWindow,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Info("gateway listening", "port", cfg.HTTPPort, "swagger", "/swagger")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down gateway...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
