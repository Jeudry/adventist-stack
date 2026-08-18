package httpx

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// In and Out carry only where a value goes, so one generic pair covers every
// operation that just moves a body. An operation that also reads the path or
// the query declares its own struct, because that is where those fields live:
// Go has no attributes, so a tag on a field is the only way to say "this comes
// from the URL".
type In[T any] struct {
	Body T
}

type Out[T any] struct {
	Body T
}

// NewAPI mounts huma on chi. huma owns binding, validation and the spec; chi
// stays the router, so the shared middleware keeps working unchanged.
func NewAPI(router *chi.Mux, title string) huma.API {
	router.Use(chimw.RequestID)
	router.Use(chimw.RealIP)
	router.Use(chimw.Recoverer)
	router.Use(chimw.Timeout(30 * time.Second))
	router.Use(Identity)

	config := huma.DefaultConfig(title, "1.0.0")
	// Drops the `$schema` field huma otherwise adds to every response body: no
	// client asked for it and the Dart models would reject it.
	config.CreateHooks = nil

	return humachi.New(router, config)
}
