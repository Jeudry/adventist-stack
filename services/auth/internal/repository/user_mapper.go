package repository

import (
	"fmt"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/auth/internal/db"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
)

func toDomain(u db.User) (domain.User, error) {
	email, err := vo.NewEmail(u.Email)
	if err != nil {
		return domain.User{}, fmt.Errorf("repository: rehydrate email: %w", err)
	}
	return domain.User{
		Base: entity.Base{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			DeletedAt: u.DeletedAt,
			CreatedBy: u.CreatedBy,
			UpdatedBy: u.UpdatedBy,
			DeletedBy: u.DeletedBy,
		},
		Email:        email,
		Name:         u.Name,
		PasswordHash: u.PasswordHash,
		Role:         domain.Role(u.Role),
	}, nil
}
