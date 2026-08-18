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

			next.ServeHTTP(w, r.WithContext(WithCaller(r.Context(), claims.Subject, claims.Role)))
		})
	}
}

// WithCaller records who the request belongs to, once the token proved it.
func WithCaller(ctx context.Context, userID, role string) context.Context {
	return context.WithValue(context.WithValue(ctx, userIDKey, userID), roleKey, role)
}

func UserID(ctx context.Context) string {
	id, _ := ctx.Value(userIDKey).(string)
	return id
}

func Role(ctx context.Context) string {
	role, _ := ctx.Value(roleKey).(string)
	return role
}

func RateLimit(requests int, window time.Duration) func(http.Handler) http.Handler {
	return httprate.LimitByIP(requests, window)
}

func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
