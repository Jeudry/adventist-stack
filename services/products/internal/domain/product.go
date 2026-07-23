package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Jeudry/adventist-stack/pkg/entity"
)

type Status string

const (
	StatusActive       Status = "active"
	StatusInactive     Status = "inactive"
	StatusDiscontinued Status = "discontinued"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusDiscontinued:
		return true
	}
	return false
}

func (s Status) String() string {
	if !s.IsValid() {
		return "unknown"
	}
	return string(s)
}

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid product")
)

const (
	NameMaxLen = 150
	SkuMaxLen  = 50
)

type Product struct {
	entity.Base
	Name        string
	Sku         string
	Description *string
	Brand       *string
	ReleaseDate *time.Time
	Status      Status
}

func (p Product) Normalize() Product {
	p.Name = strings.TrimSpace(p.Name)
	p.Sku = strings.TrimSpace(p.Sku)
	if p.Status == "" {
		p.Status = StatusActive
	}
	return p
}

func (p Product) Validate() error {
	return errors.Join(
		validateName(p.Name),
		validateSku(p.Sku),
		validateStatus(p.Status),
	)
}

func validateName(name string) error {
	switch {
	case name == "":
		return fmt.Errorf("%w: name is required", ErrInvalidProduct)
	case len(name) > NameMaxLen:
		return fmt.Errorf("%w: name exceeds %d characters", ErrInvalidProduct, NameMaxLen)
	}
	return nil
}

func validateSku(sku string) error {
	switch {
	case sku == "":
		return fmt.Errorf("%w: sku is required", ErrInvalidProduct)
	case len(sku) > SkuMaxLen:
		return fmt.Errorf("%w: sku exceeds %d characters", ErrInvalidProduct, SkuMaxLen)
	}
	return nil
}

func validateStatus(s Status) error {
	if !s.IsValid() {
		return fmt.Errorf("%w: invalid status (%q)", ErrInvalidProduct, s)
	}
	return nil
}
