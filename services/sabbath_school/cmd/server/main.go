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
	"github.com/Jeudry/adventist-stack/pkg/logger"
	sabbathschool "github.com/Jeudry/adventist-stack/services/sabbath_school"
	sabbathschoolhttp "github.com/Jeudry/adventist-stack/services/sabbath_school/internal/http"
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/repository"
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/service"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"dev"`
	HTTPPort string `env:"SABBATH_SCHOOL_HTTP_PORT" envDefault:"50056"`
	Postgres config.Postgres
}

func main() {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	log := logger.New("sabbath_school", cfg.Env)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.Postgres.AutoMigrate {
		if err := database.Migrate(cfg.Postgres.DSN, "sabbath_school_schema_migrations", sabbathschool.MigrationsFS, "migrations"); err != nil {
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

	repo := repository.NewSabbathSchoolRepository(pool)
	svc := service.NewSabbathSchoolService(repo)

	if err := httpx.Serve(ctx, cfg.HTTPPort, api(svc), log); err != nil {
		log.Error("http server", "err", err)
		os.Exit(1)
	}
}

func api(svc *service.SabbathSchoolService) *chi.Mux {
	router := chi.NewRouter()
	sabbathschoolhttp.NewHandler(svc).Register(httpx.NewAPI(router, "Sabbath School"))
	return router
}
