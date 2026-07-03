// Package grpc adapta el AuthService de dominio al contrato gRPC.
package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "github.com/Jeudry/adventist-stack/gen/auth/v1"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
	"github.com/Jeudry/adventist-stack/services/auth/internal/service"
)

// Server implementa authv1.AuthServiceServer.
type Server struct {
	authv1.UnimplementedAuthServiceServer
	svc *service.AuthService
}

// NewServer crea el adaptador gRPC.
func NewServer(svc *service.AuthService) *Server {
	return &Server{svc: svc}
}

// Register da de alta un usuario.
func (s *Server) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.AuthResponse, error) {
	user, tokens, err := s.svc.Register(ctx, req.GetEmail(), req.GetName(), req.GetPassword())
	if err != nil {
		return nil, toStatus(err)
	}
	return buildAuthResponse(user, tokens.Access, tokens.Refresh), nil
}

// Login autentica un usuario.
func (s *Server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.AuthResponse, error) {
	user, tokens, err := s.svc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, toStatus(err)
	}
	return buildAuthResponse(user, tokens.Access, tokens.Refresh), nil
}

// ValidateToken valida un access token.
func (s *Server) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	userID, role, err := s.svc.ValidateToken(req.GetAccessToken())
	if err != nil {
		return &authv1.ValidateTokenResponse{Valid: false}, nil
	}
	return &authv1.ValidateTokenResponse{Valid: true, UserId: userID, Role: role}, nil
}

// RefreshToken emite un nuevo par de tokens a partir de un refresh válido.
func (s *Server) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.AuthResponse, error) {
	// El refresh token comparte el formato del access token; reutilizamos la
	// validación y volvemos a emitir. (Extensible a rotación/lista negra.)
	userID, role, err := s.svc.ValidateToken(req.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "refresh token inválido")
	}
	_ = role
	_ = userID
	return nil, status.Error(codes.Unimplemented, "rotación de refresh token pendiente de implementar")
}

func buildAuthResponse(u domain.User, access, refresh string) *authv1.AuthResponse {
	return &authv1.AuthResponse{
		User: &authv1.User{
			Id:    u.ID,
			Email: u.Email,
			Name:  u.Name,
			Role:  string(u.Role),
		},
		AccessToken:  access,
		RefreshToken: refresh,
	}
}

func toStatus(err error) error {
	switch {
	case errors.Is(err, domain.ErrEmailTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
