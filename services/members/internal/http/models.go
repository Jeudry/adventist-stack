package http

import (
	"github.com/Jeudry/adventist-stack/pkg/httpx"
)

type MemberRequest struct {
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Gender      string  `json:"gender"`
	Address     *string `json:"address,omitempty"`
	BirthDate   *string `json:"birthDate,omitempty"`
	BaptismDate *string `json:"baptismDate,omitempty"`
	Status      string  `json:"status,omitempty"`
}

type MemberVM struct {
	httpx.BaseVM
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Gender      string  `json:"gender"`
	Address     *string `json:"address,omitempty"`
	BirthDate   *string `json:"birthDate,omitempty"`
	BaptismDate *string `json:"baptismDate,omitempty"`
	Status      string  `json:"status"`
}
