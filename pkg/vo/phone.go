package vo

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ErrInvalidPhone is the sentinel wrapped by every phone validation failure.
var ErrInvalidPhone = errors.New("invalid phone")

// Phone length bounds measured on the canonical form (optional + plus digits).
const (
	PhoneMinLen = 7
	PhoneMaxLen = 20
)

var (
	phoneRegex = regexp.MustCompile(`^\+?[0-9]+$`)
	phoneStrip = strings.NewReplacer(" ", "", "-", "", "(", "", ")", "", ".", "")
)

// Phone is a validated, canonical phone number. The zero value is invalid;
// build one only through NewPhone so the invariant always holds.
type Phone struct {
	value string
}

// NewPhone canonicalizes (strip separators) and validates a raw phone number.
func NewPhone(raw string) (Phone, error) {
	normalized := phoneStrip.Replace(strings.TrimSpace(raw))

	switch {
	case normalized == "":
		return Phone{}, fmt.Errorf("%w: phone is required", ErrInvalidPhone)
	case len(normalized) < PhoneMinLen:
		return Phone{}, fmt.Errorf("%w: phone must have at least %d digits", ErrInvalidPhone, PhoneMinLen)
	case len(normalized) > PhoneMaxLen:
		return Phone{}, fmt.Errorf("%w: phone exceeds %d characters", ErrInvalidPhone, PhoneMaxLen)
	case !phoneRegex.MatchString(normalized):
		return Phone{}, fmt.Errorf("%w: invalid format", ErrInvalidPhone)
	}

	return Phone{value: normalized}, nil
}

// NewOptionalPhone builds a Phone from an optional raw value: nil or blank
// yields the zero Phone, any other value is validated like NewPhone.
func NewOptionalPhone(raw *string) (Phone, error) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return Phone{}, nil
	}
	return NewPhone(*raw)
}

// String returns the canonical number for persistence and transport.
func (p Phone) String() string { return p.value }

// IsZero reports whether the phone is the unset zero value.
func (p Phone) IsZero() bool { return p.value == "" }

// Equals compares two phones by their canonical value.
func (p Phone) Equals(other Phone) bool { return p.value == other.value }

// Ptr returns nil for the zero Phone, otherwise a pointer to its value,
// ready to persist into a nullable column.
func (p Phone) Ptr() *string {
	if p.IsZero() {
		return nil
	}
	v := p.value
	return &v
}
