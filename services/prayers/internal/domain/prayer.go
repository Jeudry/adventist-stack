package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Jeudry/adventist-stack/pkg/entity"
)

const (
	TitleMinLen       = 3
	TitleMaxLen       = 128
	DescriptionMinLen = 5
	DescriptionMaxLen = 2056
	AuthorNameMinLen  = 3
	AuthorNameMaxLen  = 256
)

type Status int

const (
	STATUS_UNSPECIFIED Status = iota + 1
	STATUS_PENDING
	STATUS_ANSWERED
	STATUS_ARCHIVED
)

func (s Status) String() string {
	switch s {
	case STATUS_UNSPECIFIED:
		return "unspecified"
	case STATUS_PENDING:
		return "pending"
	case STATUS_ANSWERED:
		return "answered"
	case STATUS_ARCHIVED:
		return "archived"
	default:
		return "unknown"
	}
}

func (s Status) IsValid() bool {
	switch s {
	case STATUS_UNSPECIFIED, STATUS_PENDING, STATUS_ANSWERED, STATUS_ARCHIVED:
		return true
	default:
		return false
	}
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
	AuthorName  string
	IsAnonymous bool
	Status      Status
}

func (p *Prayer) Normalize() {
	p.AuthorName = strings.TrimSpace(p.AuthorName)
	p.Title = strings.TrimSpace(p.Title)
	p.Description = strings.TrimSpace(p.Description)
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
		return fmt.Errorf("title must be between %d and %d characters", TitleMinLen, TitleMaxLen)
	default:
		return nil
	}
}

func validateDescription(description string) error {
	switch {
	case len(description) < DescriptionMinLen || len(description) > DescriptionMaxLen:
		return fmt.Errorf("description must be between %d and %d characters", DescriptionMinLen, DescriptionMaxLen)
	default:
		return nil
	}
}

func validateAuthorName(authorName string) error {
	switch {
	case len(authorName) < AuthorNameMinLen || len(authorName) > AuthorNameMaxLen:
		return fmt.Errorf("author name must be between %d and %d characters", AuthorNameMinLen, AuthorNameMaxLen)
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
