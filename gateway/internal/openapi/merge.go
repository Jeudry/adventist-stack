// Package openapi builds the public API document out of the ones each service
// generates for itself.
//
// The services are the only place that knows their own shape, so the gateway
// asks them instead of keeping a copy: a hand-written document is a second
// truth, and the second truth is always the one that goes stale.
package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"sync"
	"time"
)

const (
	fetchTimeout = 5 * time.Second
	cacheFor     = 30 * time.Second
)

type Upstream struct {
	Name string
	URL  string
}

type Merger struct {
	upstreams []Upstream
	title     string
	version   string
	log       *slog.Logger
	client    *http.Client

	mu       sync.Mutex
	local    map[string]any
	cached   []byte
	cachedAt time.Time
}

func NewMerger(title, version string, upstreams []Upstream, log *slog.Logger) *Merger {
	return &Merger{
		upstreams: upstreams,
		title:     title,
		version:   version,
		log:       log,
		client:    &http.Client{Timeout: fetchTimeout},
	}
}

// Document merges every reachable service's spec. A service that is down is
// left out with a warning rather than failing the whole document: partial
// documentation beats a broken page, and this is a page humans read.
//
// The result is cached briefly so opening Swagger does not fan out on every
// asset the page requests.
func (m *Merger) Document(ctx context.Context) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cached != nil && time.Since(m.cachedAt) < cacheFor {
		return m.cached, nil
	}

	merged := map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":   m.title,
			"version": m.version,
			"description": "Generado a partir de lo que declara cada servicio. " +
				"El gateway valida el JWT, aplica el límite de peticiones y enruta.",
		},
		"servers": []any{map[string]any{"url": "/"}},
	}
	paths := map[string]any{}
	schemas := map[string]any{}

	if m.local != nil {
		collect(m.local, "paths", paths)
		if components, ok := m.local["components"].(map[string]any); ok {
			collect(components, "schemas", schemas)
		}
	}

	for _, upstream := range m.upstreams {
		spec, err := m.fetch(ctx, upstream)
		if err != nil {
			m.log.Warn("service left out of the API document", "service", upstream.Name, "err", err)
			continue
		}
		collect(spec, "paths", paths)
		if components, ok := spec["components"].(map[string]any); ok {
			collect(components, "schemas", schemas)
		}
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("openapi: no service answered")
	}

	merged["paths"] = paths
	merged["components"] = map[string]any{
		"schemas": schemas,
		"securitySchemes": map[string]any{
			"bearerAuth": map[string]any{"type": "http", "scheme": "bearer", "bearerFormat": "JWT"},
		},
	}
	merged["security"] = []any{map[string]any{"bearerAuth": []any{}}}

	document, err := json.Marshal(merged)
	if err != nil {
		return nil, fmt.Errorf("openapi: %w", err)
	}

	m.cached, m.cachedAt = document, time.Now()
	return document, nil
}

func (m *Merger) fetch(ctx context.Context, upstream Upstream) (map[string]any, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, upstream.URL+"/openapi.json", nil)
	if err != nil {
		return nil, err
	}

	response, err := m.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("answered %d", response.StatusCode)
	}

	var spec map[string]any
	if err := json.NewDecoder(response.Body).Decode(&spec); err != nil {
		return nil, err
	}
	return spec, nil
}

func collect(from map[string]any, key string, into map[string]any) {
	if found, ok := from[key].(map[string]any); ok {
		maps.Copy(into, found)
	}
}

// Handler serves the merged document.
func (m *Merger) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		document, err := m.Document(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(document)
	}
}

// IncludeLocal adds a spec the gateway builds in process — its own operations,
// which no upstream can describe because the gateway is the one serving them.
func (m *Merger) IncludeLocal(api huma.API) {
	document, err := json.Marshal(api.OpenAPI())
	if err != nil {
		m.log.Warn("could not read the gateway's own spec", "err", err)
		return
	}
	var spec map[string]any
	if err := json.Unmarshal(document, &spec); err != nil {
		m.log.Warn("could not read the gateway's own spec", "err", err)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.local = spec
	m.cached = nil
}
