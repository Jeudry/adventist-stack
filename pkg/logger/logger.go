// Package logger provee un *slog.Logger estructurado y consistente para todos
// los servicios.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New crea un logger JSON estructurado. En "dev" usa salida de texto legible.
func New(service, environment string) *slog.Logger {
	level := slog.LevelInfo
	if strings.EqualFold(environment, "dev") || strings.EqualFold(environment, "development") {
		level = slog.LevelDebug
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: level}
	if strings.HasPrefix(strings.ToLower(environment), "dev") {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler).With("service", service)
}
