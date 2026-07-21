package handlers

import (
	"net/http"

	authv1 "github.com/Jeudry/adventist-stack/gen/auth/v1"
)

type AuthHandler struct {
	auth authv1.AuthServiceClient
}

func NewAuthHandler(auth authv1.AuthServiceClient) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserVM struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type AuthResponse struct {
	User         UserVM `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
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

	writeJSON(w, http.StatusCreated, toAuthResponse(res.GetSession()))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
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

	writeJSON(w, http.StatusOK, toAuthResponse(res.GetSession()))
}

func toAuthResponse(s *authv1.Session) AuthResponse {
	return AuthResponse{
		User: UserVM{
			ID:    s.GetUser().GetId(),
			Email: s.GetUser().GetEmail(),
			Name:  s.GetUser().GetName(),
			Role:  s.GetUser().GetRole(),
		},
		AccessToken:  s.GetAccessToken(),
		RefreshToken: s.GetRefreshToken(),
	}
}
