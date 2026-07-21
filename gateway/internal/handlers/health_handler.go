package handlers

import "net/http"

func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func Me(userIDFn, roleFn func(r *http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"user_id": userIDFn(r),
			"role":    roleFn(r),
		})
	}
}
