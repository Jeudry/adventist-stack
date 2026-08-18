// Package httpx holds what every service needs to speak HTTP: JSON in, JSON out,
// and the caller the gateway already authenticated.
package httpx

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
)

// Headers the gateway asserts after verifying the token. A service behind the
// gateway trusts them and never sees the token itself, so its own port must not
// be reachable from outside the cluster.
const (
	UserIDHeader = "X-User-Id"
	RoleHeader   = "X-User-Role"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg})
}

func DecodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func ListRequest(r *http.Request) pagination.ListRequest {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	return pagination.ListRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   strings.TrimSpace(query.Get("search")),
	}
}
