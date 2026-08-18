package domain

import (
	"errors"
	"strings"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/google/uuid"
)

type Status string

const (
	NameMinLen     = 3
	NameMaxLen     = 128
	LocationMinLen = 2
	LocationMaxLen = 128
	MinAgeLimit    = 0
	MaxAgeLimit    = 120
)

const (
	STATUS_UNSPECIFIED Status = ""
	STATUS_ACTIVE      Status = "active"
	STATUS_INACTIVE    Status = "inactive"
)

func (s Status) String() string { return string(s) }

func (s Status) IsValid() bool {
	switch s {
	case STATUS_ACTIVE, STATUS_INACTIVE:
		return true
	default:
		return false
	}
}

// ParseStatus maps external text (e.g. "Active", " inactive ") to a Status, or
// the zero value when it does not match a known status.
func ParseStatus(s string) Status {
	candidate := Status(strings.ToLower(strings.TrimSpace(s)))
	if candidate.IsValid() {
		return candidate
	}
	return STATUS_UNSPECIFIED
}

var (
	ErrSabbathSchoolNotFound = errors.New("sabbath school not found")
	ErrInvalidStatus         = errors.New("invalid status")
	ErrInvalidSabbathSchool  = errors.New("invalid sabbath school")
)

type SabbathSchool struct {
	entity.Base
	Name         string
	TeacherID    uuid.NullUUID
	Location     *string
	TargetMinAge *int
	TargetMaxAge *int
	Status       Status
}

func (ss *SabbathSchool) Normalize() {
	ss.Name = strings.TrimSpace(ss.Name)
	if ss.Location != nil {
		trimmed := strings.TrimSpace(*ss.Location)
		if trimmed == "" {
			ss.Location = nil
		} else {
			ss.Location = &trimmed
		}
	}
	if ss.Status == STATUS_UNSPECIFIED {
		ss.Status = STATUS_ACTIVE
	}
}

func (ss SabbathSchool) Validate() error {
	return errors.Join(
		validateName(ss.Name),
		validateLocation(ss.Location),
		validateAgeLimits(ss.TargetMinAge, ss.TargetMaxAge),
		validateStatus(ss.Status),
	)
}

func validateName(name string) error {
	if len(name) < NameMinLen || len(name) > NameMaxLen {
		return ErrInvalidSabbathSchool
	}
	return nil
}

func validateLocation(location *string) error {
	if location == nil {
		return nil
	}
	if len(*location) < LocationMinLen || len(*location) > LocationMaxLen {
		return ErrInvalidSabbathSchool
	}
	return nil
}

func validateAgeLimits(minAge, maxAge *int) error {
	if minAge != nil && (*minAge < MinAgeLimit || *minAge > MaxAgeLimit) {
		return ErrInvalidSabbathSchool
	}
	if maxAge != nil && (*maxAge < MinAgeLimit || *maxAge > MaxAgeLimit) {
		return ErrInvalidSabbathSchool
	}
	if minAge != nil && maxAge != nil && *minAge > *maxAge {
		return ErrInvalidSabbathSchool
	}
	return nil
}

func validateStatus(status Status) error {
	if !status.IsValid() {
		return ErrInvalidStatus
	}
	return nil
}
