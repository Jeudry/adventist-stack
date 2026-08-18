package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Jeudry/adventist-stack/gateway/internal/openapi"
	"github.com/Jeudry/adventist-stack/gateway/internal/proxy"
	"github.com/Jeudry/adventist-stack/gateway/internal/router"
	"github.com/Jeudry/adventist-stack/pkg/config"
	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/jwt"
	"github.com/Jeudry/adventist-stack/pkg/logger"
)

type Config struct {
	Env              string        `env:"ENV" envDefault:"dev"`
	HTTPPort         string        `env:"GATEWAY_HTTP_PORT" envDefault:"8080"`
	AuthURL          string        `env:"AUTH_URL" envDefault:"http://localhost:50051"`
	MembersURL       string        `env:"MEMBERS_URL" envDefault:"http://localhost:50052"`
	PrayersURL       string        `env:"PRAYERS_URL" envDefault:"http://localhost:50055"`
	SabbathSchoolURL string        `env:"SABBATH_SCHOOL_URL" envDefault:"http://localhost:50056"`
	AllowedOrigins   string        `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	RateLimit        int           `env:"RATE_LIMIT_REQUESTS" envDefault:"100"`
	RateWindow       time.Duration `env:"RATE_LIMIT_WINDOW" envDefault:"1m"`
	JWT              config.JWT
}

func main() {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	log := logger.New("gateway", cfg.Env)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Every service declares its full public path, so the gateway forwards the
	// request untouched: no prefix is stripped on the way through.
	targets := map[string]string{
		"/api/v1/auth":            cfg.AuthURL,
		"/api/v1/members":         cfg.MembersURL,
		"/api/v1/prayers":         cfg.PrayersURL,
		"/api/v1/sabbath-schools": cfg.SabbathSchoolURL,
	}
	upstreams := make(map[string]http.Handler, len(targets))
	for mount, target := range targets {
		handler, err := proxy.To(target, log)
		if err != nil {
			log.Error("invalid upstream URL", "mount", mount, "target", target, "err", err)
			os.Exit(1)
		}
		upstreams[mount] = handler
	}

	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	merger := openapi.NewMerger("Adventist Stack API", "1.0.0", []openapi.Upstream{
		{Name: "auth", URL: cfg.AuthURL},
		{Name: "members", URL: cfg.MembersURL},
		{Name: "prayers", URL: cfg.PrayersURL},
		{Name: "sabbath_school", URL: cfg.SabbathSchoolURL},
	}, log)

	handler, gatewayAPI := router.New(router.Deps{
		JWT:            jwtManager,
		Auth:           upstreams["/api/v1/auth"],
		Members:        upstreams["/api/v1/members"],
		Prayers:        upstreams["/api/v1/prayers"],
		SabbathSchool:  upstreams["/api/v1/sabbath-schools"],
		OpenAPI:        merger.Handler(),
		AllowedOrigins: strings.Split(cfg.AllowedOrigins, ","),
		RateLimit:      cfg.RateLimit,
		RateWindow:     cfg.RateWindow,
	})

	merger.IncludeLocal(gatewayAPI)

	log.Info("gateway ready", "swagger", "/swagger")
	if err := httpx.Serve(ctx, cfg.HTTPPort, handler, log); err != nil {
		log.Error("http server", "err", err)
		os.Exit(1)
	}
}
