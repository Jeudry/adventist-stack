package middleware

import (
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Jeudry/adventist-stack/pkg/jwt"
)

// HumaAuth is Auth for an operation registered with huma. The chi version
// guards handlers the gateway proxies; this one guards the operations the
// gateway answers itself, and puts the caller where the handler can read it.
func HumaAuth(api huma.API, manager *jwt.Manager) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		token, ok := strings.CutPrefix(ctx.Header("Authorization"), "Bearer ")
		if !ok || token == "" {
			_ = huma.WriteErr(api, ctx, 401, "unauthorized")
			return
		}

		claims, err := manager.Verify(token)
		if err != nil {
			_ = huma.WriteErr(api, ctx, 401, "unauthorized")
			return
		}

		next(huma.WithValue(huma.WithValue(ctx, userIDKey, claims.Subject), roleKey, claims.Role))
	}
}
