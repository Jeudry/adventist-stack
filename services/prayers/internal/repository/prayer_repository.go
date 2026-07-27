package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/db"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PrayerRepository struct {
	db *db.Queries
}

func NewPrayerRepository(pool *pgxpool.Pool) *PrayerRepository {
	return &PrayerRepository{db: db.New(pool)}
}

func (r *PrayerRepository) Create(ctx context.Context, prayer domain.Prayer) (domain.Prayer, error) {
	created, err := r.db.CreatePayer(ctx, toCreateParams(prayer))
	if err != nil {
		return domain.Prayer{}, fmt.Errorf("repository: failed to create prayer: %w", err)
	}
	return toDomain(created)
}

func (r *PrayerRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Prayer, error) {
	p, err := r.db.GetPrayerByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Prayer{}, domain.ErrPrayerNotFound
	}
	if err != nil {
		return domain.Prayer{}, fmt.Errorf("repository: failed to get prayer by id: %w", err)
	}
	return toDomain(p)
}

func (r *PrayerRepository) RetrieveList(ctx context.Context, pq pagination.Query) ([]domain.Prayer, error) {
	pqParams := toListParams(pq)
	rows, err := r.db.ListPrayers(ctx, pqParams)
	if err != nil {
		return nil, fmt.Errorf("repostiory: failed to list prayers: %w", err)
	}

	out := make([]domain.Prayer, len(rows))
	for i, row := range rows {
		prayerRow, err := toDomain(row)
		if err != nil {
			return nil, err
		}
		out[i] = prayerRow
	}
	return out, nil
}

func (r *PrayerRepository) Count(ctx context.Context, pq pagination.Query) (int, error) {
	count, err := r.db.Count(ctx, pq.Search)
	if err != nil {
		return 0, fmt.Errorf("repository: failed to count prayers: %w", err)
	}
	return int(count), nil
}

func (r *PrayerRepository) Update(ctx context.Context, prayer domain.Prayer) (domain.Prayer, error) {
	toUpdateParams := toUpdateParams(prayer)
	updated, err := r.db.UpdatePrayer(ctx, toUpdateParams)

	if err != nil {
		return domain.Prayer{}, fmt.Errorf("repository: failed to update prayer: %w", err)
	}
	return toDomain(updated)
}

func (r *PrayerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.DeletePrayer(ctx, id)
	if err != nil {
		return fmt.Errorf("repository: failed to delete prayer: %w", err)
	}
	return nil
}
