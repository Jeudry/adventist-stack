package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusVisitor  Status = "visitor"
)

var ErrMemberNotFound = errors.New("miembro no encontrado")

type Member struct {
	Id          uuid.UUID
	FirstName   string
	LastName    string
	Email       *string
	Phone       *string
	Gender      string
	Address     *string
	BirthDate   *time.Time
	BaptismDate *time.Time
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
