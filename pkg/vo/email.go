// Package vo holds reusable value objects shared across services.
package vo

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ErrInvalidEmail is the sentinel wrapped by every email validation failure.
var ErrInvalidEmail = errors.New("invalid email")

// EmailMaxLen is the maximum total length allowed by RFC 5321.
const EmailMaxLen = 254

var emailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)

// Email is a validated, normalized email address. The zero value is invalid;
// build one only through NewEmail so the invariant always holds.
type Email struct {
	value string
}

// NewEmail normalizes (trim + lowercase) and validates a raw email.
func NewEmail(raw string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))

	switch {
	case normalized == "":
		return Email{}, fmt.Errorf("%w: email is required", ErrInvalidEmail)
	case len(normalized) > EmailMaxLen:
		return Email{}, fmt.Errorf("%w: email exceeds %d characters", ErrInvalidEmail, EmailMaxLen)
	case !emailRegex.MatchString(normalized):
		return Email{}, fmt.Errorf("%w: invalid format", ErrInvalidEmail)
	}

	return Email{value: normalized}, nil
}

// String returns the normalized address for persistence and transport.
func (e Email) String() string { return e.value }

// IsZero reports whether the email is the unset zero value.
func (e Email) IsZero() bool { return e.value == "" }

// Equals compares two emails by their normalized value.
func (e Email) Equals(other Email) bool { return e.value == other.value }
