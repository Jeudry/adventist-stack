package httpx

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

const shutdownGrace = 15 * time.Second

// Serve runs handler until ctx is cancelled, then drains in-flight requests.
// It returns once the server has stopped, so callers can defer their cleanup.
func Serve(ctx context.Context, port string, handler http.Handler, log *slog.Logger) error {
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errs := make(chan error, 1)
	go func() {
		log.Info("listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- err
		}
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
	}

	log.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGrace)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
