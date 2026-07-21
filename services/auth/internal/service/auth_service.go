package service

import (
	"context"

	"github.com/Jeudry/adventist-stack/pkg/jwt"
	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
)

type userRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, email vo.Email) (domain.User, error)
	ExistsByEmail(ctx context.Context, email vo.Email) (bool, error)
}

type AuthService struct {
	repo userRepository
	jwt  *jwt.Manager
}

func New(repo userRepository, jwtManager *jwt.Manager) *AuthService {
	return &AuthService{repo: repo, jwt: jwtManager}
}

type Tokens struct {
	Access  string
	Refresh string
}

func (s *AuthService) Register(ctx context.Context, email, name, password string) (domain.User, Tokens, error) {
	user, err := domain.NewUser(email, name, domain.RoleMember)
	if err != nil {
		return domain.User{}, Tokens{}, err
	}
	if err := user.SetPassword(password); err != nil {
		return domain.User{}, Tokens{}, err
	}

	exists, err := s.repo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return domain.User{}, Tokens{}, err
	}
	if exists {
		return domain.User{}, Tokens{}, domain.ErrEmailTaken
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return domain.User{}, Tokens{}, err
	}

	tokens, err := s.issue(created)
	return created, tokens, err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (domain.User, Tokens, error) {
	parsed, err := vo.NewEmail(email)
	if err != nil {
		return domain.User{}, Tokens{}, domain.ErrInvalidCredentials
	}

	user, err := s.repo.FindByEmail(ctx, parsed)
	if err != nil {
		return domain.User{}, Tokens{}, domain.ErrInvalidCredentials
	}

	if err := user.Authenticate(password); err != nil {
		return domain.User{}, Tokens{}, domain.ErrInvalidCredentials
	}

	tokens, err := s.issue(user)
	return user, tokens, err
}

func (s *AuthService) ValidateToken(accessToken string) (userID, role string, err error) {
	claims, err := s.jwt.Verify(accessToken)
	if err != nil {
		return "", "", err
	}
	return claims.Subject, claims.Role, nil
}

func (s *AuthService) issue(u domain.User) (Tokens, error) {
	access, refresh, err := s.jwt.Generate(u.ID.String(), string(u.Role))
	if err != nil {
		return Tokens{}, err
	}
	return Tokens{Access: access, Refresh: refresh}, nil
}
