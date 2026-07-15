// Package service contiene la lógica de negocio del servicio auth.
package service

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/Jeudry/adventist-stack/pkg/jwt"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
)

// userRepository es la dependencia de persistencia (interface para testear).
type userRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// AuthService orquesta registro, login y emisión de tokens.
type AuthService struct {
	repo userRepository
	jwt  *jwt.Manager
}

// New crea el servicio.
func New(repo userRepository, jwtManager *jwt.Manager) *AuthService {
	return &AuthService{repo: repo, jwt: jwtManager}
}

// Tokens agrupa el par de tokens emitido.
type Tokens struct {
	Access  string
	Refresh string
}

// Register crea un usuario nuevo con la contraseña hasheada y devuelve tokens.
func (s *AuthService) Register(ctx context.Context, email, name, password string) (domain.User, Tokens, error) {
	email = normalizeEmail(email)
	if email == "" || strings.TrimSpace(name) == "" || len(password) < 8 {
		return domain.User{}, Tokens{}, fmt.Errorf("datos inválidos: email, nombre y contraseña (>=8) son obligatorios")
	}

	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return domain.User{}, Tokens{}, err
	}
	if exists {
		return domain.User{}, Tokens{}, domain.ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, Tokens{}, fmt.Errorf("service: hash password: %w", err)
	}

	created, err := s.repo.Create(ctx, domain.User{
		Email:        email,
		Name:         strings.TrimSpace(name),
		PasswordHash: string(hash),
		Role:         domain.RoleMember,
	})
	if err != nil {
		return domain.User{}, Tokens{}, err
	}

	tokens, err := s.issue(created)
	return created, tokens, err
}

// Login valida credenciales y devuelve tokens.
func (s *AuthService) Login(ctx context.Context, email, password string) (domain.User, Tokens, error) {
	email = normalizeEmail(email)

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		// No revelamos si el usuario existe o no.
		return domain.User{}, Tokens{}, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.User{}, Tokens{}, domain.ErrInvalidCredentials
	}

	tokens, err := s.issue(user)
	return user, tokens, err
}

// ValidateToken verifica un access token y devuelve el userID y el rol.
func (s *AuthService) ValidateToken(accessToken string) (userID, role string, err error) {
	claims, err := s.jwt.Verify(accessToken)
	if err != nil {
		return "", "", err
	}
	return claims.Subject, claims.Role, nil
}

func (s *AuthService) issue(u domain.User) (Tokens, error) {
	// El JWT usa el id como string (el "subject" del token).
	access, refresh, err := s.jwt.Generate(u.ID.String(), string(u.Role))
	if err != nil {
		return Tokens{}, err
	}
	return Tokens{Access: access, Refresh: refresh}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
