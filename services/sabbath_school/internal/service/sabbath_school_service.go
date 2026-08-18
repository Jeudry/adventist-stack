package service

import (
	"context"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/domain"
	"github.com/google/uuid"
)

type SabbathSchoolRepository interface {
	Create(context.Context, domain.SabbathSchool) (domain.SabbathSchool, error)
	GetByID(context.Context, uuid.UUID) (domain.SabbathSchool, error)
	RetrieveList(context.Context, pagination.Query) ([]domain.SabbathSchool, error)
	Count(context.Context, pagination.Query) (int, error)
	Update(context.Context, domain.SabbathSchool) (domain.SabbathSchool, error)
	Delete(ctx context.Context, id, deletedBy uuid.UUID) error
}

type SabbathSchoolService struct {
	repo SabbathSchoolRepository
}

func NewSabbathSchoolService(repo SabbathSchoolRepository) *SabbathSchoolService {
	return &SabbathSchoolService{
		repo: repo,
	}
}

func (s *SabbathSchoolService) Create(ctx context.Context, ss domain.SabbathSchool) (domain.SabbathSchool, error) {
	ss.Normalize()
	if err := ss.Validate(); err != nil {
		return domain.SabbathSchool{}, err
	}
	return s.repo.Create(ctx, ss)
}

func (s *SabbathSchoolService) GetByID(ctx context.Context, id uuid.UUID) (domain.SabbathSchool, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SabbathSchoolService) RetrieveList(ctx context.Context, pq pagination.Query) (pagination.Page[domain.SabbathSchool], error) {
	items, err := s.repo.RetrieveList(ctx, pq)
	if err != nil {
		return pagination.Page[domain.SabbathSchool]{}, err
	}
	total, err := s.repo.Count(ctx, pq)
	if err != nil {
		return pagination.Page[domain.SabbathSchool]{}, err
	}

	return pagination.NewPage(items, total, pq), nil
}

func (s *SabbathSchoolService) Update(ctx context.Context, sc domain.SabbathSchool) (domain.SabbathSchool, error) {
	if sc.ID == uuid.Nil {
		return domain.SabbathSchool{}, domain.ErrSabbathSchoolNotFound
	}

	sc.Normalize()
	if err := sc.Validate(); err != nil {
		return domain.SabbathSchool{}, err
	}

	return s.repo.Update(ctx, sc)
}

func (s *SabbathSchoolService) Delete(ctx context.Context, id, deletedBy uuid.UUID) error {
	return s.repo.Delete(ctx, id, deletedBy)
}
