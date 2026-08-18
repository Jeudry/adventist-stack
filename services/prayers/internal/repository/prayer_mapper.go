package repository

import (
	"github.com/jackc/pgx/v5"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/domain"
)

// row is what pgx.Row and pgx.Rows have in common, so one scan serves the single-row queries and
// the list alike.
type row interface {
	Scan(dest ...any) error
}

// The order here must match prayerColumns exactly. A NULL author_name is an anonymous prayer, and
// it stays a nil pointer all the way into the entity.
func scanPrayer(r row) (domain.Prayer, error) {
	var (
		p         domain.Prayer
		base      entity.Base
		rawStatus string
	)

	err := r.Scan(
		&base.ID, &p.Title, &p.Description, &p.AuthorName, &rawStatus,
		&base.CreatedAt, &base.CreatedBy, &base.UpdatedAt, &base.UpdatedBy, &base.DeletedAt, &base.DeletedBy,
	)
	if err != nil {
		return domain.Prayer{}, err
	}

	p.Base = base
	p.Status = domain.Status(rawStatus)
	return p, nil
}

// writableArgs carries the columns the client owns. The insert takes it whole; the update swaps
// created_by for updated_by, so neither statement can quietly write the other's audit column.
func writableArgs(p domain.Prayer) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"title":       p.Title,
		"description": p.Description,
		"author_name": p.AuthorName,
		"status":      p.Status.String(),
		"created_by":  p.CreatedBy,
	}
}
