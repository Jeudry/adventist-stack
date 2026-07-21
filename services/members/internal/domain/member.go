package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/strutil"
	"github.com/Jeudry/adventist-stack/pkg/vo"
)

const (
	NamesMaxLen     = 256
	NamesMinLen     = 1
	AddressMaxLen   = 1024
	AddressMinLen   = 5
	MaxBirthdateAge = 140
	MaxBaptismAge   = 130
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusVisitor  Status = "visitor"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusVisitor:
		return true
	default:
		return false
	}
}

type Gender string

const (
	GenderMale   Gender = "M"
	GenderFemale Gender = "F"
)

func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale:
		return true
	default:
		return false
	}
}

var (
	ErrMemberNotFound  = errors.New("Member not found")
	ErrorInvalidMember = errors.New("Member invalid")
)

type Member struct {
	entity.Base
	FirstName   string
	LastName    string
	Email       vo.Email
	Phone       vo.Phone
	Gender      Gender
	Address     *string
	BirthDate   *time.Time
	BaptismDate *time.Time
	Status      Status
}

func (m Member) Normalize() Member {
	m.FirstName = strings.TrimSpace(m.FirstName)
	m.LastName = strings.TrimSpace(m.LastName)
	m.Gender = Gender(strings.ToUpper(strings.TrimSpace(string(m.Gender))))
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
		validateGender(m.Gender),
		validateAddress(m.Address),
		validateBirthDate(m.BirthDate),
		validateBaptismDate(m.BaptismDate),
		validateStatus(m.Status),
	)
}

func validateBaptismDate(baptismDate *time.Time) error {
	switch {
	case baptismDate == nil:
		return nil
	case baptismDate.After(time.Now()):
		return fmt.Errorf("%w: baptism date cannot be in the future", ErrorInvalidMember)
	case baptismDate.Before(time.Now().AddDate(-MaxBaptismAge, 0, 0)):
		return fmt.Errorf("%w: baptism date cannot be more than %d years ago", ErrorInvalidMember, MaxBaptismAge)
	default:
		return nil
	}
}

func validateBirthDate(birthDate *time.Time) error {
	switch {
	case birthDate == nil:
		return nil
	case birthDate.After(time.Now()):
		return fmt.Errorf("%w: baptism date cannot be in the future", ErrorInvalidMember)
	case birthDate.Before(time.Now().AddDate(-MaxBirthdateAge, 0, 0)):
		return fmt.Errorf("%w: baptism date cannot be more than %d years ago", ErrorInvalidMember, MaxBirthdateAge)
	default:
		return nil
	}
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

func validateGender(gender Gender) error {
	if !gender.IsValid() {
		return fmt.Errorf("%w: invalid gender", ErrorInvalidMember)
	}
	return nil
}

func validateAddress(address *string) error {
	switch {
	case len(*address) < AddressMinLen:
		return fmt.Errorf("%w: address must be at least %d characters", ErrorInvalidMember, AddressMinLen)
	case len(*address) > AddressMaxLen:
		return fmt.Errorf("%w: address must be at most %d characters", ErrorInvalidMember, AddressMaxLen)
	}
	return nil
}

func validateStatus(status Status) error {
	if !status.IsValid() {
		return fmt.Errorf("%w: invalid status", ErrorInvalidMember)
	}
	return nil
}
