// Package models contiene los modelos de persistencia (mapeo directo a las
// columnas de la base de datos) del servicio auth.
package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
)

// User mapea la tabla users.
type User struct {
	ID           uuid.UUID
	Email        string
	Name         string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ToDomain convierte el modelo de DB al modelo de dominio.
func (u User) ToDomain() domain.User {
	return domain.User{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		PasswordHash: u.PasswordHash,
		Role:         domain.Role(u.Role),
	}
}
