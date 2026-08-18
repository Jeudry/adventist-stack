package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/services/auth/internal/service"
)

type Handler struct {
	svc *service.AuthService
}

func NewHandler(svc *service.AuthService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "register",
		Method:        "POST",
		Path:          "/api/v1/auth/register",
		Summary:       "Registrar un usuario",
		Tags:          []string{"auth"},
		DefaultStatus: 201,
		Security:      []map[string][]string{},
	}, h.SignUp)

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      "POST",
		Path:        "/api/v1/auth/login",
		Summary:     "Iniciar sesión",
		Tags:        []string{"auth"},
		Security:    []map[string][]string{},
	}, h.SignIn)
}

func (h *Handler) SignUp(ctx context.Context, in *httpx.In[RegisterRequest]) (*httpx.Out[AuthResponse], error) {
	user, tokens, err := h.svc.Register(ctx, in.Body.Email, in.Body.Name, in.Body.Password)
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[AuthResponse]{Body: toAuthResponse(user, tokens)}, nil
}

func (h *Handler) SignIn(ctx context.Context, in *httpx.In[LoginRequest]) (*httpx.Out[AuthResponse], error) {
	user, tokens, err := h.svc.Login(ctx, in.Body.Email, in.Body.Password)
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[AuthResponse]{Body: toAuthResponse(user, tokens)}, nil
}
