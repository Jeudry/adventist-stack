// Package middleware agrupa los middlewares HTTP reutilizables del gateway:
// autenticación JWT, rate limiting y CORS.
package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"

	"github.com/Jeudry/adventist-stack/pkg/jwt"
)

type ctxKey string

const (
	userIDKey ctxKey = "userID"
	roleKey   ctxKey = "role"
)

// Auth valida el header Authorization: Bearer <token> y guarda el userID y
// el role en el contexto. Rechaza con 401 si el token es inválido o falta.
func Auth(manager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := manager.Verify(token)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.Subject)
			ctx = context.WithValue(ctx, roleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserID extrae el ID del usuario autenticado del contexto.
func UserID(ctx context.Context) string {
	id, _ := ctx.Value(userIDKey).(string)
	return id
}

// Role extrae el rol del usuario autenticado del contexto.
func Role(ctx context.Context) string {
	role, _ := ctx.Value(roleKey).(string)
	return role
}

// RateLimit limita las peticiones por IP (requests por ventana de tiempo).
func RateLimit(requests int, window time.Duration) func(http.Handler) http.Handler {
	return httprate.LimitByIP(requests, window)
}

// CORS devuelve un middleware CORS configurado a partir de los orígenes dados.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
