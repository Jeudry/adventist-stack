package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Jeudry/adventist-stack/pkg/vo"
)

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidUser        = errors.New("invalid user")
)

const (
	NameMaxLen     = 150
	PasswordMinLen = 8
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleMember:
		return true
	default:
		return false
	}
}

type User struct {
	ID           uuid.UUID
	Email        vo.Email
	Name         string
	PasswordHash string
	Role         Role
}

// NewUser builds a valid User or returns why it could not.
func NewUser(rawEmail, name string, role Role) (User, error) {
	email, err := vo.NewEmail(rawEmail)
	if err != nil {
		return User{}, err
	}
	name = strings.TrimSpace(name)

	switch {
	case name == "":
		return User{}, fmt.Errorf("%w: name is required", ErrInvalidUser)
	case len(name) > NameMaxLen:
		return User{}, fmt.Errorf("%w: name exceeds %d characters", ErrInvalidUser, NameMaxLen)
	case !role.IsValid():
		return User{}, fmt.Errorf("%w: invalid role (%q)", ErrInvalidUser, role)
	}

	return User{Email: email, Name: name, Role: role}, nil
}

// SetPassword enforces the password policy and stores its hash.
func (u *User) SetPassword(plain string) error {
	if len(plain) < PasswordMinLen {
		return fmt.Errorf("%w: password must have at least %d characters", ErrInvalidUser, PasswordMinLen)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	u.PasswordHash = string(hash)
	return nil
}

// Authenticate reports whether the plain password matches the stored hash.
func (u User) Authenticate(plain string) error {
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plain)) != nil {
		return ErrInvalidCredentials
	}
	return nil
}

// IsAdmin reports whether the user holds the admin role.
func (u User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
