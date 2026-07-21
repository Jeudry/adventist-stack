package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "github.com/Jeudry/adventist-stack/gen/auth/v1"
	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
	"github.com/Jeudry/adventist-stack/services/auth/internal/service"
)

type Server struct {
	authv1.UnimplementedAuthServiceServer
	svc *service.AuthService
}

func NewServer(svc *service.AuthService) *Server {
	return &Server{svc: svc}
}

func (s *Server) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user, tokens, err := s.svc.Register(ctx, req.GetEmail(), req.GetName(), req.GetPassword())
	if err != nil {
		return nil, toStatus(err)
	}
	return &authv1.RegisterResponse{Session: buildSession(user, tokens.Access, tokens.Refresh)}, nil
}

func (s *Server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	user, tokens, err := s.svc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, toStatus(err)
	}
	return &authv1.LoginResponse{Session: buildSession(user, tokens.Access, tokens.Refresh)}, nil
}

func (s *Server) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	userID, role, err := s.svc.ValidateToken(req.GetAccessToken())
	if err != nil {
		return &authv1.ValidateTokenResponse{Valid: false}, nil
	}
	return &authv1.ValidateTokenResponse{Valid: true, UserId: userID, Role: role}, nil
}

func (s *Server) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	userID, role, err := s.svc.ValidateToken(req.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}
	_ = role
	_ = userID
	return nil, status.Error(codes.Unimplemented, "refresh token rotation not yet implemented")
}

func buildSession(u domain.User, access, refresh string) *authv1.Session {
	return &authv1.Session{
		User: &authv1.User{
			Id:    u.ID.String(),
			Email: u.Email.String(),
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
	case errors.Is(err, domain.ErrInvalidUser), errors.Is(err, vo.ErrInvalidEmail):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
