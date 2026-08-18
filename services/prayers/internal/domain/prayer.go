package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/strutil"
)

const (
	TitleMinLen       = 3
	TitleMaxLen       = 128
	DescriptionMinLen = 5
	DescriptionMaxLen = 2056
	AuthorNameMinLen  = 3
	AuthorNameMaxLen  = 256
)

type Status string

const (
	STATUS_UNSPECIFIED Status = ""
	STATUS_PENDING     Status = "pending"
	STATUS_ANSWERED    Status = "answered"
	STATUS_ARCHIVED    Status = "archived"
)

func (s Status) String() string { return string(s) }

func (s Status) IsValid() bool {
	switch s {
	case STATUS_PENDING, STATUS_ANSWERED, STATUS_ARCHIVED:
		return true
	default:
		return false
	}
}

// ParseStatus maps external text (e.g. "Pending", " archived ") to a Status, or
// the zero value when it does not match a known status.
func ParseStatus(s string) Status {
	candidate := Status(strings.ToLower(strings.TrimSpace(s)))
	if candidate.IsValid() {
		return candidate
	}
	return STATUS_UNSPECIFIED
}

var (
	ErrPrayerNotFound = errors.New("prayer not found")
	ErrInvalidStatus  = errors.New("invalid status")
	ErrInvalidPrayer  = errors.New("invalid prayer")
)

type Prayer struct {
	entity.Base
	Title       string
	Description string
	// nil is an anonymous prayer. One field instead of a name plus a flag, so
	// "anonymous but with a name" is not a state anything can construct.
	AuthorName *string
	Status     Status
}

func (p *Prayer) Normalize() {
	// A name that trims to nothing is the same as not having one.
	p.AuthorName = strutil.TrimPtr(p.AuthorName)
	p.Title = strings.TrimSpace(p.Title)
	p.Description = strings.TrimSpace(p.Description)

	// A request that omits the status gets the one it starts life with, the
	// same value the column defaults to. Without this the API rejects a body
	// the spec declares valid.
	if p.Status == STATUS_UNSPECIFIED {
		p.Status = STATUS_PENDING
	}
}

func (p Prayer) Validate() error {
	return errors.Join(
		validateTitle(p.Title),
		validateDescription(p.Description),
		validateAuthorName(p.AuthorName),
		validateStatus(p.Status),
	)
}

func validateTitle(title string) error {
	switch {
	case len(title) < TitleMinLen || len(title) > TitleMaxLen:
		return fmt.Errorf("%w: title must be between %d and %d characters", ErrInvalidPrayer, TitleMinLen, TitleMaxLen)
	default:
		return nil
	}
}

func validateDescription(description string) error {
	switch {
	case len(description) < DescriptionMinLen || len(description) > DescriptionMaxLen:
		return fmt.Errorf("%w: description must be between %d and %d characters", ErrInvalidPrayer, DescriptionMinLen, DescriptionMaxLen)
	default:
		return nil
	}
}

// No author at all is valid — that is an anonymous prayer, and it keeps no name
// to leak. A name that is present has to be a real one.
func validateAuthorName(authorName *string) error {
	switch {
	case authorName == nil:
		return nil
	case len(*authorName) < AuthorNameMinLen || len(*authorName) > AuthorNameMaxLen:
		return fmt.Errorf("%w: author name must be between %d and %d characters", ErrInvalidPrayer, AuthorNameMinLen, AuthorNameMaxLen)
	default:
		return nil
	}
}

func validateStatus(status Status) error {
	switch {
	case !status.IsValid():
		return ErrInvalidStatus
	default:
		return nil
	}
}
