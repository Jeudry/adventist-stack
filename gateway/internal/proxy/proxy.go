// Package proxy forwards a request to a service without reading its body.
// The gateway decides who the caller is; the payload is none of its business.
package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/middleware"
)

// To returns a handler that forwards the request to target untouched, asserting
// the caller in headers the service trusts. The path is not rewritten: every
// service declares the same public path the client called.
//
// Whatever identity headers arrived from outside are dropped first: they are
// the gateway's word, and a client must not be able to speak in its place.
func To(target string, log *slog.Logger) (http.Handler, error) {
	parsed, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	reverse := httputil.NewSingleHostReverseProxy(parsed)
	reverse.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Error("upstream unreachable", "target", target, "path", r.URL.Path, "err", err)
		httpx.WriteError(w, http.StatusBadGateway, "service unavailable")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del(httpx.UserIDHeader)
		r.Header.Del(httpx.RoleHeader)
		if userID := middleware.UserID(r.Context()); userID != "" {
			r.Header.Set(httpx.UserIDHeader, userID)
			r.Header.Set(httpx.RoleHeader, middleware.Role(r.Context()))
		}
		reverse.ServeHTTP(w, r)
	}), nil
}
