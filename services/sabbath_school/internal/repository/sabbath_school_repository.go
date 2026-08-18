package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/pkg/ptr"
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/domain"
)

// Every statement names its parameters and pgx.StrictNamedArgs refuses to run when one is missing
// or spelled wrong. Positional placeholders fail silently instead: target_min_age and
// target_max_age are both integers, and swapping them still runs.
const (
	sabbathSchoolColumns = `id, name, teacher_id, location, target_min_age, target_max_age, status,
	                        created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`

	sabbathSchoolSearch = `(
	@search = ''
	OR name ILIKE ('%' || @search || '%')
	OR location ILIKE ('%' || @search || '%')
	OR status ILIKE ('%' || @search || '%')
) AND deleted_at IS NULL`

	insertSabbathSchool = `
INSERT INTO sabbath_school (name, teacher_id, location, target_min_age, target_max_age, status, created_by)
VALUES (@name, @teacher_id, @location, @target_min_age, @target_max_age, @status, @created_by)
RETURNING ` + sabbathSchoolColumns

	selectSabbathSchool = `
SELECT ` + sabbathSchoolColumns + `
FROM sabbath_school
WHERE id = @id AND deleted_at IS NULL`

	listSabbathSchools = `
SELECT ` + sabbathSchoolColumns + `
FROM sabbath_school
WHERE ` + sabbathSchoolSearch + `
ORDER BY created_at DESC
LIMIT @row_limit OFFSET @row_offset`

	countSabbathSchools = `SELECT count(*) FROM sabbath_school WHERE ` + sabbathSchoolSearch

	updateSabbathSchool = `
UPDATE sabbath_school SET
    name           = @name,
    teacher_id     = @teacher_id,
    location       = @location,
    target_min_age = @target_min_age,
    target_max_age = @target_max_age,
    status         = @status,
    updated_at     = NOW(),
    updated_by     = @updated_by
WHERE id = @id AND deleted_at IS NULL
RETURNING ` + sabbathSchoolColumns

	deleteSabbathSchool = `
UPDATE sabbath_school SET deleted_at = NOW(), deleted_by = @deleted_by
WHERE id = @id AND deleted_at IS NULL`

	// Every other query filters deleted_at IS NULL, which is what makes a soft delete a delete.
	// This one does not, so a test can check what the repository buried.
	selectSabbathSchoolIncludingDeleted = `
SELECT ` + sabbathSchoolColumns + ` FROM sabbath_school WHERE id = @id`
)

type SabbathSchoolRepository struct {
	pool *pgxpool.Pool
}

func NewSabbathSchoolRepository(pool *pgxpool.Pool) *SabbathSchoolRepository {
	return &SabbathSchoolRepository{pool: pool}
}

func (r *SabbathSchoolRepository) Create(ctx context.Context, ss domain.SabbathSchool) (domain.SabbathSchool, error) {
	created, err := scanSabbathSchool(r.pool.QueryRow(ctx, insertSabbathSchool, writableArgs(ss)))
	if err != nil {
		return domain.SabbathSchool{}, fmt.Errorf("repository: create sabbath school: %w", err)
	}
	return created, nil
}

func (r *SabbathSchoolRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.SabbathSchool, error) {
	found, err := scanSabbathSchool(r.pool.QueryRow(ctx, selectSabbathSchool, pgx.StrictNamedArgs{"id": id}))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.SabbathSchool{}, domain.ErrSabbathSchoolNotFound
	}
	if err != nil {
		return domain.SabbathSchool{}, fmt.Errorf("repository: get sabbath school: %w", err)
	}
	return found, nil
}

func (r *SabbathSchoolRepository) RetrieveList(ctx context.Context, q pagination.Query) ([]domain.SabbathSchool, error) {
	rows, err := r.pool.Query(ctx, listSabbathSchools, pgx.StrictNamedArgs{
		"search":     q.Search,
		"row_limit":  q.Limit,
		"row_offset": q.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("repository: list sabbath schools: %w", err)
	}
	defer rows.Close()

	schools := make([]domain.SabbathSchool, 0)
	for rows.Next() {
		school, err := scanSabbathSchool(rows)
		if err != nil {
			return nil, fmt.Errorf("repository: list sabbath schools: %w", err)
		}
		schools = append(schools, school)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: list sabbath schools: %w", err)
	}
	return schools, nil
}

func (r *SabbathSchoolRepository) Count(ctx context.Context, q pagination.Query) (int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, countSabbathSchools, pgx.StrictNamedArgs{"search": q.Search}).Scan(&total); err != nil {
		return 0, fmt.Errorf("repository: count sabbath schools: %w", err)
	}
	return total, nil
}

func (r *SabbathSchoolRepository) Update(ctx context.Context, ss domain.SabbathSchool) (domain.SabbathSchool, error) {
	args := writableArgs(ss)
	delete(args, "created_by")
	args["id"] = ss.ID
	args["updated_by"] = ss.UpdatedBy

	updated, err := scanSabbathSchool(r.pool.QueryRow(ctx, updateSabbathSchool, args))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.SabbathSchool{}, domain.ErrSabbathSchoolNotFound
	}
	if err != nil {
		return domain.SabbathSchool{}, fmt.Errorf("repository: update sabbath school: %w", err)
	}
	return updated, nil
}

func (r *SabbathSchoolRepository) Delete(ctx context.Context, id, deletedBy uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, deleteSabbathSchool, pgx.StrictNamedArgs{
		"id":         id,
		"deleted_by": ptr.NonZero(deletedBy),
	})
	if err != nil {
		return fmt.Errorf("repository: delete sabbath school: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrSabbathSchoolNotFound
	}
	return nil
}
