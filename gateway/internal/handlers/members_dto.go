package handlers

type CreateMemberRequest struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Gender      string  `json:"gender"`
	Address     *string `json:"address,omitempty"`
	BirthDate   *string `json:"birth_date,omitempty"`   // Formato "YYYY-MM-DD"
	BaptismDate *string `json:"baptism_date,omitempty"` // Formato "YYYY-MM-DD"
	Status      string  `json:"status,omitempty"`
}

type UpdateMemberRequest struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Gender      string  `json:"gender"`
	Address     *string `json:"address,omitempty"`
	BirthDate   *string `json:"birth_date,omitempty"`   // Formato "YYYY-MM-DD"
	BaptismDate *string `json:"baptism_date,omitempty"` // Formato "YYYY-MM-DD"
	Status      string  `json:"status,omitempty"`
}

type MemberVM struct {
	BaseVM
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
