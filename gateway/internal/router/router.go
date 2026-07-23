package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Jeudry/adventist-stack/gateway/api"
	"github.com/Jeudry/adventist-stack/gateway/internal/handlers"
	"github.com/Jeudry/adventist-stack/pkg/jwt"
	"github.com/Jeudry/adventist-stack/pkg/middleware"
)

type Deps struct {
	JWT            *jwt.Manager
	AuthHandler    *handlers.AuthHandler
	MembersHandler *handlers.MembersHandler
	AllowedOrigins []string
	RateLimit      int
	RateWindow     time.Duration
}

func New(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(middleware.CORS(d.AllowedOrigins))
	r.Use(middleware.RateLimit(d.RateLimit, d.RateWindow))

	r.Get("/health", handlers.Health)

	r.Get("/swagger", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(api.SwaggerHTML)
	})
	r.Get("/openapi.yaml", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(api.OpenAPISpec)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", d.AuthHandler.Register)
		r.Post("/auth/login", d.AuthHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(d.JWT))
			r.Get("/me", handlers.Me(
				func(req *http.Request) string { return middleware.UserID(req.Context()) },
				func(req *http.Request) string { return middleware.Role(req.Context()) },
			))

			r.Route("/members", func(r chi.Router) {
				r.Post("/", d.MembersHandler.Create)
				r.Get("/", d.MembersHandler.List)
				r.Get("/{id}", d.MembersHandler.GetByID)
				r.Put("/{id}", d.MembersHandler.Update)
				r.Delete("/{id}", d.MembersHandler.Delete)
			})
		})
	})

	return r
}
