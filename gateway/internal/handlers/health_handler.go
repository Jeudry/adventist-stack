package handlers

import "net/http"

// Health maneja GET /health (liveness probe).
func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Me maneja GET /api/v1/me y devuelve el usuario del token (ruta protegida).
// Extrae el userID/role del contexto poblado por el middleware de auth.
func Me(userIDFn, roleFn func(r *http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"user_id": userIDFn(r),
			"role":    roleFn(r),
		})
	}
}
