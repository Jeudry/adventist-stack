package repository

import (
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
)

// row is what pgx.Row and pgx.Rows have in common, so one scan serves the single-row queries and
// the list alike.
type row interface {
	Scan(dest ...any) error
}

// The order here must match memberColumns exactly. Email and phone are value objects with their
// own validation, which is why the crossing can fail on the way out of the database.
func scanMember(r row) (domain.Member, error) {
	var (
		m         domain.Member
		base      entity.Base
		rawEmail  *string
		rawPhone  *string
		rawGender string
		rawStatus string
	)

	err := r.Scan(
		&base.ID, &m.FirstName, &m.LastName, &rawEmail, &rawPhone, &rawGender, &rawStatus,
		&m.BirthDate, &m.BaptismDate, &m.Address,
		&base.CreatedAt, &base.CreatedBy, &base.UpdatedAt, &base.UpdatedBy, &base.DeletedAt, &base.DeletedBy,
	)
	if err != nil {
		return domain.Member{}, err
	}

	email, err := vo.NewOptionalEmail(rawEmail)
	if err != nil {
		return domain.Member{}, fmt.Errorf("rehydrate email: %w", err)
	}
	phone, err := vo.NewOptionalPhone(rawPhone)
	if err != nil {
		return domain.Member{}, fmt.Errorf("rehydrate phone: %w", err)
	}

	m.Base = base
	m.Email = email
	m.Phone = phone
	m.Gender = domain.Gender(rawGender)
	m.Status = domain.Status(rawStatus)
	return m, nil
}

// writableArgs carries the columns the client owns. The insert takes it whole; the update swaps
// created_by for updated_by, so neither statement can quietly write the other's audit column.
func writableArgs(m domain.Member) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"first_name":   m.FirstName,
		"last_name":    m.LastName,
		"email":        m.Email.Ptr(),
		"phone":        m.Phone.Ptr(),
		"gender":       m.Gender.String(),
		"address":      m.Address,
		"birth_date":   m.BirthDate,
		"baptism_date": m.BaptismDate,
		"status":       m.Status.String(),
		"created_by":   m.CreatedBy,
	}
}
