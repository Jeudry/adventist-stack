package handlers

import (
	"net/http"

	"github.com/Jeudry/adventist-stack/pkg/httpx"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func Me(userIDFn, roleFn func(r *http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{
			"userId": userIDFn(r),
			"role":   roleFn(r),
		})
	}
}
