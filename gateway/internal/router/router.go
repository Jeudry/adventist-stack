// Package router construye el router HTTP del gateway: middlewares globales,
// rutas públicas y protegidas, y Swagger UI.
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

// Deps son las dependencias que el router necesita.
type Deps struct {
	JWT            *jwt.Manager
	AuthHandler    *handlers.AuthHandler
	AllowedOrigins []string
	RateLimit      int
	RateWindow     time.Duration
}

// New arma el http.Handler completo.
func New(d Deps) http.Handler {
	r := chi.NewRouter()

	// Middlewares globales.
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(middleware.CORS(d.AllowedOrigins))
	r.Use(middleware.RateLimit(d.RateLimit, d.RateWindow))

	// Salud.
	r.Get("/health", handlers.Health)

	// Swagger UI + spec.
	r.Get("/swagger", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(api.SwaggerHTML)
	})
	r.Get("/openapi.yaml", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(api.OpenAPISpec)
	})

	// API v1.
	r.Route("/api/v1", func(r chi.Router) {
		// Público.
		r.Post("/auth/register", d.AuthHandler.Register)
		r.Post("/auth/login", d.AuthHandler.Login)

		// Protegido (requiere JWT válido).
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(d.JWT))
			r.Get("/me", handlers.Me(
				func(req *http.Request) string { return middleware.UserID(req.Context()) },
				func(req *http.Request) string { return middleware.Role(req.Context()) },
			))
		})
	})

	return r
}
