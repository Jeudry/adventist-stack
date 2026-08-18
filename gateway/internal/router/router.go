package router

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Jeudry/adventist-stack/gateway/api"
	"github.com/Jeudry/adventist-stack/gateway/internal/handlers"
	"github.com/Jeudry/adventist-stack/pkg/jwt"
	"github.com/Jeudry/adventist-stack/pkg/middleware"
)

type Deps struct {
	JWT            *jwt.Manager
	Auth           http.Handler
	Members        http.Handler
	Prayers        http.Handler
	SabbathSchool  http.Handler
	OpenAPI        http.HandlerFunc
	AllowedOrigins []string
	RateLimit      int
	RateWindow     time.Duration
}

// New returns the gateway handler plus the spec of what the gateway answers
// itself, so it can be merged with the ones the services generate.
func New(d Deps) (http.Handler, huma.API) {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(middleware.CORS(d.AllowedOrigins))
	r.Use(middleware.RateLimit(d.RateLimit, d.RateWindow))

	r.Get("/swagger", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(api.SwaggerHTML)
	})
	r.Get("/openapi.json", d.OpenAPI)

	// The gateway's own two operations. This API only describes them; the
	// document clients read is the merged one served above.
	config := huma.DefaultConfig("Gateway", "1.0.0")
	config.CreateHooks = nil
	config.DocsPath = ""
	config.OpenAPIPath = ""
	gatewayAPI := humachi.New(r, config)
	handlers.Register(gatewayAPI, d.JWT)

	// Registration and login are the only way in, so they are the only routes
	// reachable without a token.
	r.Mount("/api/v1/auth", d.Auth)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(d.JWT))
		r.Mount("/api/v1/members", d.Members)
		r.Mount("/api/v1/prayers", d.Prayers)
		r.Mount("/api/v1/sabbath-schools", d.SabbathSchool)
	})

	return r, gatewayAPI
}
