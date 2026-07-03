// Package domain contiene los modelos de core (dominio) del servicio auth,
// independientes de la base de datos y del transporte.
package domain

import "errors"

// Errores de dominio reutilizables por las capas superiores.
var (
	ErrEmailTaken         = errors.New("el email ya está registrado")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUserNotFound       = errors.New("usuario no encontrado")
)

// Role define los roles de usuario del sistema.
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

// User es el modelo de dominio. No conoce columnas ni JSON.
type User struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
	Role         Role
}
