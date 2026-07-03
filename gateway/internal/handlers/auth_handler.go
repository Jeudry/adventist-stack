package handlers

import (
	"net/http"

	authv1 "github.com/Jeudry/adventist-stack/gen/auth/v1"
)

// AuthHandler expone los endpoints REST de autenticación, delegando en el
// microservicio auth vía gRPC.
type AuthHandler struct {
	auth authv1.AuthServiceClient
}

// NewAuthHandler crea el handler.
func NewAuthHandler(auth authv1.AuthServiceClient) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// --- Viewmodels (DTOs) ---

// RegisterRequest es el cuerpo de POST /api/v1/auth/register.
type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// LoginRequest es el cuerpo de POST /api/v1/auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserVM es la representación pública de un usuario.
type UserVM struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// AuthResponse es la respuesta de register/login.
type AuthResponse struct {
	User         UserVM `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Register maneja POST /api/v1/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "JSON inválido"})
		return
	}

	res, err := h.auth.Register(r.Context(), &authv1.RegisterRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toAuthResponse(res))
}

// Login maneja POST /api/v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "JSON inválido"})
		return
	}

	res, err := h.auth.Login(r.Context(), &authv1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toAuthResponse(res))
}

func toAuthResponse(res *authv1.AuthResponse) AuthResponse {
	return AuthResponse{
		User: UserVM{
			ID:    res.GetUser().GetId(),
			Email: res.GetUser().GetEmail(),
			Name:  res.GetUser().GetName(),
			Role:  res.GetUser().GetRole(),
		},
		AccessToken:  res.GetAccessToken(),
		RefreshToken: res.GetRefreshToken(),
	}
}
