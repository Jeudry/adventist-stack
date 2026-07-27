package member

import "github.com/Jeudry/adventist-stack/gateway/internal/models/base"

type MemberVM struct {
	base.BaseVM
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Gender      string  `json:"gender"`
	Address     *string `json:"address,omitempty"`
	BirthDate   *string `json:"birth_date,omitempty"`   // Formato "YYYY-MM-DD"
	BaptismDate *string `json:"baptism_date,omitempty"` // Formato "YYYY-MM-DD"
	Status      string  `json:"status"`
}
