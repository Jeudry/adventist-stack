package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/products/internal/domain"
)

type productRepository interface {
	Create(ctx context.Context, p domain.Product) (domain.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
	List(ctx context.Context, q pagination.Query) ([]domain.Product, error)
	Count(ctx context.Context, q pagination.Query) (int, error)
	Update(ctx context.Context, p domain.Product) (domain.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProductService struct {
	repo productRepository
}

func New(repo productRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	p = p.Normalize()
	if err := p.Validate(); err != nil {
		return domain.Product{}, err
	}
	return s.repo.Create(ctx, p)
}

func (s *ProductService) Get(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context, req pagination.ListRequest) (pagination.Page[domain.Product], error) {
	q := req.ToQuery()

	items, err := s.repo.List(ctx, q)
	if err != nil {
		return pagination.Page[domain.Product]{}, err
	}
	total, err := s.repo.Count(ctx, q)
	if err != nil {
		return pagination.Page[domain.Product]{}, err
	}

	return pagination.NewPage(items, total, q), nil
}

func (s *ProductService) Update(ctx context.Context, p domain.Product) (domain.Product, error) {
	if p.Id == uuid.Nil {
		return domain.Product{}, fmt.Errorf("%w: id is required", domain.ErrInvalidProduct)
	}
	p = p.Normalize()
	if err := p.Validate(); err != nil {
		return domain.Product{}, err
	}
	return s.repo.Update(ctx, p)
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
