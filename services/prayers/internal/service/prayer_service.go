package service

import (
	"context"
	"fmt"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/domain"
	"github.com/google/uuid"
)

type prayerRepository interface {
	Create(ctx context.Context, p domain.Prayer) (domain.Prayer, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Prayer, error)
	RetrieveList(ctx context.Context, pq pagination.Query) ([]domain.Prayer, error)
	Count(ctx context.Context, pq pagination.Query) (int, error)
	Update(ctx context.Context, p domain.Prayer) (domain.Prayer, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type PrayerService struct {
	prayerRepo prayerRepository
}

func NewPrayerService(prayerRepo prayerRepository) *PrayerService {
	return &PrayerService{
		prayerRepo: prayerRepo,
	}
}

func (s *PrayerService) Create(ctx context.Context, p domain.Prayer) (domain.Prayer, error) {
	p.Normalize()
	if err := p.Validate(); err != nil {
		return domain.Prayer{}, err
	}
	return s.prayerRepo.Create(ctx, p)
}

func (s *PrayerService) GetByID(ctx context.Context, id uuid.UUID) (domain.Prayer, error) {
	return s.prayerRepo.GetByID(ctx, id)
}

func (s *PrayerService) RetrieveList(ctx context.Context, q pagination.Query) (pagination.Page[domain.Prayer], error) {
	items, err := s.prayerRepo.RetrieveList(ctx, q)
	if err != nil {
		return pagination.Page[domain.Prayer]{}, err
	}
	total, err := s.prayerRepo.Count(ctx, q)
	if err != nil {
		return pagination.Page[domain.Prayer]{}, err
	}

	return pagination.NewPage(items, total, q), nil
}
func (s *PrayerService) Update(ctx context.Context, p domain.Prayer) (domain.Prayer, error) {
	if p.ID == uuid.Nil {
		return domain.Prayer{}, fmt.Errorf("%w: id is required for update", domain.ErrPrayerNotFound)
	}

	p.Normalize()
	if err := p.Validate(); err != nil {
		return domain.Prayer{}, err
	}

	return s.prayerRepo.Update(ctx, p)
}

func (s *PrayerService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.prayerRepo.Delete(ctx, id)
}
