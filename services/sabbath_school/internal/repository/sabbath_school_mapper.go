package repository

import (
	"github.com/jackc/pgx/v5"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/domain"
)

// row is what pgx.Row and pgx.Rows have in common, so one scan serves the single-row queries and
// the list alike.
type row interface {
	Scan(dest ...any) error
}

// The order here must match sabbathSchoolColumns exactly. teacher_id is a uuid.NullUUID because a
// class without a teacher is a real state, not a missing value.
func scanSabbathSchool(r row) (domain.SabbathSchool, error) {
	var (
		ss        domain.SabbathSchool
		base      entity.Base
		rawStatus string
	)

	err := r.Scan(
		&base.ID, &ss.Name, &ss.TeacherID, &ss.Location, &ss.TargetMinAge, &ss.TargetMaxAge, &rawStatus,
		&base.CreatedAt, &base.CreatedBy, &base.UpdatedAt, &base.UpdatedBy, &base.DeletedAt, &base.DeletedBy,
	)
	if err != nil {
		return domain.SabbathSchool{}, err
	}

	ss.Base = base
	ss.Status = domain.Status(rawStatus)
	return ss, nil
}

// writableArgs carries the columns the client owns. The insert takes it whole; the update swaps
// created_by for updated_by, so neither statement can quietly write the other's audit column.
func writableArgs(ss domain.SabbathSchool) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"name":           ss.Name,
		"teacher_id":     ss.TeacherID,
		"location":       ss.Location,
		"target_min_age": ss.TargetMinAge,
		"target_max_age": ss.TargetMaxAge,
		"status":         ss.Status.String(),
		"created_by":     ss.CreatedBy,
	}
}
