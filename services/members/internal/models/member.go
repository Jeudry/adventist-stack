package models

import (
	"time"

	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
	"github.com/google/uuid"
)

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
	Status      domain.Status
	CreateAt    time.Time
	UpdatedAt   time.Time
}
