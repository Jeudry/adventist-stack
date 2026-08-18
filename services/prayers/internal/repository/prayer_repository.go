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
	"github.com/Jeudry/adventist-stack/services/prayers/internal/domain"
)

// Every statement names its parameters and pgx.StrictNamedArgs refuses to run when one is missing
// or spelled wrong. Positional placeholders fail silently instead: title and description are both
// text, and swapping them still runs.
const (
	prayerColumns = `id, title, description, author_name, status,
	                 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`

	prayerSearch = `(
	@search = ''
	OR title ILIKE ('%' || @search || '%')
	OR description ILIKE ('%' || @search || '%')
	OR author_name ILIKE ('%' || @search || '%')
) AND deleted_at IS NULL`

	insertPrayer = `
INSERT INTO prayers (title, description, author_name, status, created_by)
VALUES (@title, @description, @author_name, @status, @created_by)
RETURNING ` + prayerColumns

	selectPrayer = `
SELECT ` + prayerColumns + `
FROM prayers
WHERE id = @id AND deleted_at IS NULL`

	listPrayers = `
SELECT ` + prayerColumns + `
FROM prayers
WHERE ` + prayerSearch + `
ORDER BY created_at DESC
LIMIT @row_limit OFFSET @row_offset`

	countPrayers = `SELECT count(*) FROM prayers WHERE ` + prayerSearch

	updatePrayer = `
UPDATE prayers SET
    title       = @title,
    description = @description,
    author_name = @author_name,
    status      = @status,
    updated_at  = NOW(),
    updated_by  = @updated_by
WHERE id = @id AND deleted_at IS NULL
RETURNING ` + prayerColumns

	deletePrayer = `
UPDATE prayers SET deleted_at = NOW(), deleted_by = @deleted_by
WHERE id = @id AND deleted_at IS NULL`

	// Every other query filters deleted_at IS NULL, which is what makes a soft delete a delete.
	// This one does not, so a test can check what the repository buried.
	selectPrayerIncludingDeleted = `
SELECT ` + prayerColumns + ` FROM prayers WHERE id = @id`
)

type PrayerRepository struct {
	pool *pgxpool.Pool
}

func NewPrayerRepository(pool *pgxpool.Pool) *PrayerRepository {
	return &PrayerRepository{pool: pool}
}

func (r *PrayerRepository) Create(ctx context.Context, p domain.Prayer) (domain.Prayer, error) {
	created, err := scanPrayer(r.pool.QueryRow(ctx, insertPrayer, writableArgs(p)))
	if err != nil {
		return domain.Prayer{}, fmt.Errorf("repository: create prayer: %w", err)
	}
	return created, nil
}

func (r *PrayerRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Prayer, error) {
	found, err := scanPrayer(r.pool.QueryRow(ctx, selectPrayer, pgx.StrictNamedArgs{"id": id}))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Prayer{}, domain.ErrPrayerNotFound
	}
	if err != nil {
		return domain.Prayer{}, fmt.Errorf("repository: get prayer: %w", err)
	}
	return found, nil
}

func (r *PrayerRepository) RetrieveList(ctx context.Context, q pagination.Query) ([]domain.Prayer, error) {
	rows, err := r.pool.Query(ctx, listPrayers, pgx.StrictNamedArgs{
		"search":     q.Search,
		"row_limit":  q.Limit,
		"row_offset": q.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("repository: list prayers: %w", err)
	}
	defer rows.Close()

	prayers := make([]domain.Prayer, 0)
	for rows.Next() {
		prayer, err := scanPrayer(rows)
		if err != nil {
			return nil, fmt.Errorf("repository: list prayers: %w", err)
		}
		prayers = append(prayers, prayer)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: list prayers: %w", err)
	}
	return prayers, nil
}

func (r *PrayerRepository) Count(ctx context.Context, q pagination.Query) (int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, countPrayers, pgx.StrictNamedArgs{"search": q.Search}).Scan(&total); err != nil {
		return 0, fmt.Errorf("repository: count prayers: %w", err)
	}
	return total, nil
}

func (r *PrayerRepository) Update(ctx context.Context, p domain.Prayer) (domain.Prayer, error) {
	args := writableArgs(p)
	delete(args, "created_by")
	args["id"] = p.ID
	args["updated_by"] = p.UpdatedBy

	updated, err := scanPrayer(r.pool.QueryRow(ctx, updatePrayer, args))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Prayer{}, domain.ErrPrayerNotFound
	}
	if err != nil {
		return domain.Prayer{}, fmt.Errorf("repository: update prayer: %w", err)
	}
	return updated, nil
}

func (r *PrayerRepository) Delete(ctx context.Context, id, deletedBy uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, deletePrayer, pgx.StrictNamedArgs{
		"id":         id,
		"deleted_by": ptr.NonZero(deletedBy),
	})
	if err != nil {
		return fmt.Errorf("repository: delete prayer: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPrayerNotFound
	}
	return nil
}
