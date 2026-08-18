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
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
)

// Every statement names its parameters and pgx.StrictNamedArgs refuses to run when one is missing
// or spelled wrong. Positional placeholders fail silently instead: first_name and last_name are
// both text, and swapping them still runs.
const (
	memberColumns = `id, first_name, last_name, email, phone, gender, status,
	                 birth_date, baptism_date, address,
	                 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`

	memberSearch = `(
	@search = ''
	OR first_name ILIKE ('%' || @search || '%')
	OR last_name ILIKE ('%' || @search || '%')
	OR email ILIKE ('%' || @search || '%')
) AND deleted_at IS NULL`

	insertMember = `
INSERT INTO members (first_name, last_name, email, phone, gender, address,
                     birth_date, baptism_date, status, created_by)
VALUES (@first_name, @last_name, @email, @phone, @gender, @address,
        @birth_date, @baptism_date, @status, @created_by)
RETURNING ` + memberColumns

	selectMember = `
SELECT ` + memberColumns + `
FROM members
WHERE id = @id AND deleted_at IS NULL`

	listMembers = `
SELECT ` + memberColumns + `
FROM members
WHERE ` + memberSearch + `
ORDER BY created_at DESC
LIMIT @row_limit OFFSET @row_offset`

	countMembers = `SELECT count(*) FROM members WHERE ` + memberSearch

	updateMember = `
UPDATE members SET
    first_name   = @first_name,
    last_name    = @last_name,
    email        = @email,
    phone        = @phone,
    gender       = @gender,
    address      = @address,
    birth_date   = @birth_date,
    baptism_date = @baptism_date,
    status       = @status,
    updated_at   = NOW(),
    updated_by   = @updated_by
WHERE id = @id AND deleted_at IS NULL
RETURNING ` + memberColumns

	deleteMember = `
UPDATE members SET deleted_at = NOW(), deleted_by = @deleted_by
WHERE id = @id AND deleted_at IS NULL`

	// Every other query filters deleted_at IS NULL, which is what makes a soft delete a delete.
	// This one does not, so a test can check what the repository buried.
	selectMemberIncludingDeleted = `
SELECT ` + memberColumns + ` FROM members WHERE id = @id`
)

type MemberRepository struct {
	pool *pgxpool.Pool
}

func NewMemberRepository(pool *pgxpool.Pool) *MemberRepository {
	return &MemberRepository{pool: pool}
}

func (r *MemberRepository) Create(ctx context.Context, m domain.Member) (domain.Member, error) {
	created, err := scanMember(r.pool.QueryRow(ctx, insertMember, writableArgs(m)))
	if err != nil {
		return domain.Member{}, fmt.Errorf("repository: create member: %w", err)
	}
	return created, nil
}

func (r *MemberRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Member, error) {
	found, err := scanMember(r.pool.QueryRow(ctx, selectMember, pgx.StrictNamedArgs{"id": id}))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Member{}, domain.ErrMemberNotFound
	}
	if err != nil {
		return domain.Member{}, fmt.Errorf("repository: get member: %w", err)
	}
	return found, nil
}

func (r *MemberRepository) RetrieveList(ctx context.Context, q pagination.Query) ([]domain.Member, error) {
	rows, err := r.pool.Query(ctx, listMembers, pgx.StrictNamedArgs{
		"search":     q.Search,
		"row_limit":  q.Limit,
		"row_offset": q.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("repository: list members: %w", err)
	}
	defer rows.Close()

	members := make([]domain.Member, 0)
	for rows.Next() {
		member, err := scanMember(rows)
		if err != nil {
			return nil, fmt.Errorf("repository: list members: %w", err)
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: list members: %w", err)
	}
	return members, nil
}

func (r *MemberRepository) Count(ctx context.Context, q pagination.Query) (int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, countMembers, pgx.StrictNamedArgs{"search": q.Search}).Scan(&total); err != nil {
		return 0, fmt.Errorf("repository: count members: %w", err)
	}
	return total, nil
}

func (r *MemberRepository) Update(ctx context.Context, m domain.Member) (domain.Member, error) {
	args := writableArgs(m)
	delete(args, "created_by")
	args["id"] = m.ID
	args["updated_by"] = m.UpdatedBy

	updated, err := scanMember(r.pool.QueryRow(ctx, updateMember, args))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Member{}, domain.ErrMemberNotFound
	}
	if err != nil {
		return domain.Member{}, fmt.Errorf("repository: update member: %w", err)
	}
	return updated, nil
}

func (r *MemberRepository) Delete(ctx context.Context, id, deletedBy uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, deleteMember, pgx.StrictNamedArgs{
		"id":         id,
		"deleted_by": ptr.NonZero(deletedBy),
	})
	if err != nil {
		return fmt.Errorf("repository: delete member: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrMemberNotFound
	}
	return nil
}
