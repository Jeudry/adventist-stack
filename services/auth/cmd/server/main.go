package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/Jeudry/adventist-stack/pkg/config"
	"github.com/Jeudry/adventist-stack/pkg/database"
	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/jwt"
	"github.com/Jeudry/adventist-stack/pkg/logger"
	auth "github.com/Jeudry/adventist-stack/services/auth"
	authhttp "github.com/Jeudry/adventist-stack/services/auth/internal/http"
	"github.com/Jeudry/adventist-stack/services/auth/internal/repository"
	"github.com/Jeudry/adventist-stack/services/auth/internal/service"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"dev"`
	HTTPPort string `env:"AUTH_HTTP_PORT" envDefault:"50051"`
	Postgres config.Postgres
	JWT      config.JWT
}

func main() {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	log := logger.New("auth", cfg.Env)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.Postgres.AutoMigrate {
		if err := database.Migrate(cfg.Postgres.DSN, "auth_schema_migrations", auth.MigrationsFS, "migrations"); err != nil {
			log.Error("migrations failed", "err", err)
			os.Exit(1)
		}
		log.Info("migrations applied")
	}

	pool, err := database.Connect(ctx, cfg.Postgres.DSN)
	if err != nil {
		log.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	repo := repository.NewUserRepository(pool)
	svc := service.New(repo, jwtManager)

	if err := httpx.Serve(ctx, cfg.HTTPPort, api(svc), log); err != nil {
		log.Error("http server", "err", err)
		os.Exit(1)
	}
}

func api(svc *service.AuthService) *chi.Mux {
	router := chi.NewRouter()
	authhttp.NewHandler(svc).Register(httpx.NewAPI(router, "Auth"))
	return router
}
