package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"html/template"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jeudry/adventist-stack/pkg/config"
	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/logger"
	"github.com/Jeudry/adventist-stack/pkg/mailer"
	"github.com/Jeudry/adventist-stack/pkg/redis"
	notifications "github.com/Jeudry/adventist-stack/services/notifications"
	notifhttp "github.com/Jeudry/adventist-stack/services/notifications/internal/http"
	"github.com/Jeudry/adventist-stack/services/notifications/internal/service"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"dev"`
	HTTPPort string `env:"NOTIFICATIONS_HTTP_PORT" envDefault:"50054"`
	Redis    config.Redis
	SMTP     config.SMTP
}

func main() {
	cfg, err := config.Load[Config]()
	if err != nil {
		panic(err)
	}

	log := logger.New("notifications", cfg.Env)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	rdb, err := redis.Connect(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Error("failed to connect to redis", "err", err)
		os.Exit(1)
	}
	defer rdb.Close()

	tmpl, err := template.ParseFS(notifications.TemplatesFS, "templates/*.html")
	if err != nil {
		log.Error("failed to parse templates", "err", err)
		os.Exit(1)
	}
	mail := mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.From, tmpl)

	svc := service.New(mail, rdb)

	if err := httpx.Serve(ctx, cfg.HTTPPort, api(svc), log); err != nil {
		log.Error("http server", "err", err)
		os.Exit(1)
	}
}

func api(svc *service.NotificationService) *chi.Mux {
	router := chi.NewRouter()
	notifhttp.NewHandler(svc).Register(httpx.NewAPI(router, "Notifications"))
	return router
}
