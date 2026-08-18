package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jeudry/adventist-stack/pkg/config"
	"github.com/Jeudry/adventist-stack/pkg/database"
	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/logger"
	members "github.com/Jeudry/adventist-stack/services/members"
	membershttp "github.com/Jeudry/adventist-stack/services/members/internal/http"
	"github.com/Jeudry/adventist-stack/services/members/internal/repository"
	"github.com/Jeudry/adventist-stack/services/members/internal/service"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"dev"`
	HTTPPort string `env:"MEMBERS_HTTP_PORT" envDefault:"50052"`
	Postgres config.Postgres
}

func main() {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	log := logger.New("members", cfg.Env)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.Postgres.AutoMigrate {
		if err := database.Migrate(cfg.Postgres.DSN, "members_schema_migrations", members.MigrationsFS, "migrations"); err != nil {
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

	repo := repository.NewMemberRepository(pool)
	svc := service.NewMemberService(repo)

	if err := httpx.Serve(ctx, cfg.HTTPPort, api(svc), log); err != nil {
		log.Error("http server", "err", err)
		os.Exit(1)
	}
}

func api(svc *service.MemberService) *chi.Mux {
	router := chi.NewRouter()
	membershttp.NewHandler(svc).Register(httpx.NewAPI(router, "Members"))
	return router
}
