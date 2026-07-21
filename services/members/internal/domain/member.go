package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Jeudry/adventist-stack/pkg/strutil"
	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/google/uuid"
)

const (
	NamesMaxLen     = 256
	NamesMinLen     = 1
	PhoneMaxLen     = 20
	PhoneMinLen     = 7
	GenderMaxLen    = 1
	GenderMinLen    = 1
	AddressMaxLen   = 1024
	AddressMinLen   = 5
	MaxBirthdateAge = 120
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusVisitor  Status = "visitor"
)

var (
	ErrMemberNotFound  = errors.New("Member not found")
	ErrorInvalidMember = errors.New("Member invalid")
)

type Member struct {
	Id          uuid.UUID
	FirstName   string
	LastName    string
	Email       vo.Email
	Phone       *string
	Gender      string
	Address     *string
	BirthDate   *time.Time
	BaptismDate *time.Time
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m Member) Normalize() Member {
	m.FirstName = strings.TrimSpace(m.FirstName)
	m.LastName = strings.TrimSpace(m.LastName)
	m.Phone = strutil.TrimPtr(m.Phone)
	m.Gender = strings.TrimSpace(m.Gender)
	m.Address = strutil.TrimPtr(m.Address)

	if m.Status == "" {
		m.Status = StatusActive
	}

	return m
}

func (m Member) Validate() error {
	return errors.Join(
		validateFirstName(m.FirstName),
		validateLastName(m.LastName),
		validatePhone(m.Phone),
		validateGender(m.Gender),
		validateAddress(m.Address),
		validateStatus(m.Status),
	)
}

func validateFirstName(firstName string) error {
	switch {
	case len(firstName) < NamesMinLen:
		return fmt.Errorf("%w: first name must be at least %d characters", ErrorInvalidMember, NamesMinLen)
	case len(firstName) > NamesMaxLen:
		return fmt.Errorf("%w: first name must be at least %d characters", ErrorInvalidMember, NamesMaxLen)
	}

	return nil
}

func validateLastName(lastName string) error {
	switch {
	case len(lastName) < NamesMinLen:
		return fmt.Errorf("%w: last name must be at least %d characters", ErrorInvalidMember, NamesMinLen)
	case len(lastName) > NamesMaxLen:
		return fmt.Errorf("%w: last name must be at least %d characters", ErrorInvalidMember, NamesMaxLen)
	}

	return nil
}

func validatePhone(phone *string) error {

}

func validateGender(gender string) error {

}

func validateAddress(address *string) error {

}

func validateStatus(status Status) error {

}
