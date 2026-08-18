package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
	"github.com/Jeudry/adventist-stack/services/auth/internal/service"
)

func toAuthResponse(u domain.User, t service.Tokens) AuthResponse {
	return AuthResponse{
		User: UserVM{
			ID:    u.ID.String(),
			Email: u.Email.String(),
			Name:  u.Name,
			Role:  u.Role.String(),
		},
		AccessToken:  t.Access,
		RefreshToken: t.Refresh,
	}
}

func domainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrEmailTaken):
		return huma.Error409Conflict(err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		return huma.Error401Unauthorized(err.Error())
	case errors.Is(err, domain.ErrUserNotFound):
		return huma.Error404NotFound(err.Error())
	case errors.Is(err, domain.ErrInvalidUser), errors.Is(err, vo.ErrInvalidEmail):
		return huma.Error422UnprocessableEntity(err.Error())
	default:
		return huma.Error500InternalServerError("internal error")
	}
}
