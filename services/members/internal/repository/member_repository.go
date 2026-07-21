package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/members/internal/db"
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
)

type MemberRepository struct {
	q *db.Queries
}

func NewMemberRepository(pool *pgxpool.Pool) *MemberRepository {
	return &MemberRepository{q: db.New(pool)}
}

func (r *MemberRepository) Create(ctx context.Context, m domain.Member) (domain.Member, error) {
	created, err := r.q.CreateMember(ctx, toCreateParams(m))
	if err != nil {
		return domain.Member{}, fmt.Errorf("repository: create member: %w", err)
	}

	return toDomain(created), nil
}

func (r *MemberRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Member, error) {
	p, err := r.q.GetMemberByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Member{}, domain.ErrMemberNotFound
	}
	if err != nil {
		return domain.Member{}, fmt.Errorf("repository: get member: %w", err)
	}

	return toDomain(p), nil
}

func (r *MemberRepository) RetrieveList(ctx context.Context, q pagination.Query) ([]domain.Member, error) {
	rows, err := r.q.ListMembers(ctx, toMemberListParams(q))
	if err != nil {
		return nil, fmt.Errorf("repository: list members: %w", err)
	}

	out := make([]domain.Member, len(rows))
	for i, row := range rows {
		out[i] = toDomain(row)
	}
	return out, nil
}

func (r *MemberRepository) Count(ctx context.Context, q pagination.Query) (int, error) {
	count, err := r.q.CountMembers(ctx, q.Search)
	if err != nil {
		return 0, fmt.Errorf("repository: count members: %w", err)
	}

	return int(count), nil
}

func (r *MemberRepository) Update(ctx context.Context, m domain.Member) (domain.Member, error) {
	updated, err := r.q.UpdateMember(ctx, toUpdateParams(m))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Member{}, domain.ErrMemberNotFound
	}
	if err != nil {
		return domain.Member{}, fmt.Errorf("repository: update member: %w", err)
	}
	return toDomain(updated), nil
}

func (r *MemberRepository) Delete(ctx context.Context, id uuid.UUID) error {
	affected, err := r.q.DeleteMember(ctx, id)
	if err != nil {
		return fmt.Errorf("repository: delete member: %w", err)
	}
	if affected == 0 {
		return domain.ErrMemberNotFound
	}
	return nil
}
