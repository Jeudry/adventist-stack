package repository

import (
	"fmt"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
)

// row is what pgx.Row and pgx.Rows have in common, so one scan serves both the single-row queries
// and any future list.
type row interface {
	Scan(dest ...any) error
}

// The order here must match userColumns exactly. Email is a value object with its own validation,
// which is why the crossing can fail on the way out of the database and not only on the way in.
func scanUser(r row) (domain.User, error) {
	var (
		u        domain.User
		rawEmail string
		rawRole  string
		base     entity.Base
	)

	err := r.Scan(
		&base.ID, &rawEmail, &u.Name, &u.PasswordHash, &rawRole,
		&base.CreatedAt, &base.CreatedBy, &base.UpdatedAt, &base.UpdatedBy, &base.DeletedAt, &base.DeletedBy,
	)
	if err != nil {
		return domain.User{}, err
	}

	email, err := vo.NewEmail(rawEmail)
	if err != nil {
		return domain.User{}, fmt.Errorf("rehydrate email: %w", err)
	}

	u.Base = base
	u.Email = email
	u.Role = domain.Role(rawRole)
	return u, nil
}
