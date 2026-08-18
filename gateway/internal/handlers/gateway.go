package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/jwt"
	"github.com/Jeudry/adventist-stack/pkg/middleware"
)

// The two endpoints the gateway answers itself. They go through huma like any
// service, so they end up in the generated document instead of being a special
// case somebody has to remember to write down.

type healthResponse struct {
	Status string `json:"status" example:"ok"`
}

type identityResponse struct {
	UserID string `json:"userId" format:"uuid"`
	Role   string `json:"role" enum:"admin,member"`
}

func Register(api huma.API, manager *jwt.Manager) {
	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      "GET",
		Path:        "/health",
		Summary:     "Liveness probe",
		Tags:        []string{"gateway"},
		Security:    []map[string][]string{},
	}, func(context.Context, *struct{}) (*httpx.Out[healthResponse], error) {
		return &httpx.Out[healthResponse]{Body: healthResponse{Status: "ok"}}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "me",
		Method:      "GET",
		Path:        "/api/v1/me",
		Summary:     "Identidad del token presentado",
		Tags:        []string{"gateway"},
		Middlewares: huma.Middlewares{middleware.HumaAuth(api, manager)},
	}, func(ctx context.Context, _ *struct{}) (*httpx.Out[identityResponse], error) {
		return &httpx.Out[identityResponse]{Body: identityResponse{
			UserID: middleware.UserID(ctx),
			Role:   middleware.Role(ctx),
		}}, nil
	})
}
