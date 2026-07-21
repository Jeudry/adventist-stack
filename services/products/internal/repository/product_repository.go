package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/products/internal/db"
	"github.com/Jeudry/adventist-stack/services/products/internal/domain"
)

type ProductRepository struct {
	q *db.Queries
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{q: db.New(pool)}
}

func (r *ProductRepository) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	created, err := r.q.CreateProduct(ctx, toCreateParams(p))
	if err != nil {
		return domain.Product{}, fmt.Errorf("repository: create product: %w", err)
	}
	return toDomain(created), nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	p, err := r.q.GetProductByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, domain.ErrProductNotFound
	}
	if err != nil {
		return domain.Product{}, fmt.Errorf("repository: get product: %w", err)
	}
	return toDomain(p), nil
}

func (r *ProductRepository) List(ctx context.Context, q pagination.Query) ([]domain.Product, error) {
	rows, err := r.q.ListProducts(ctx, toListParams(q))
	if err != nil {
		return nil, fmt.Errorf("repository: list products: %w", err)
	}
	out := make([]domain.Product, len(rows))
	for i, p := range rows {
		out[i] = toDomain(p)
	}
	return out, nil
}

func (r *ProductRepository) Count(ctx context.Context, q pagination.Query) (int, error) {
	n, err := r.q.CountProducts(ctx, q.Search)
	if err != nil {
		return 0, fmt.Errorf("repository: count products: %w", err)
	}
	return int(n), nil
}

func (r *ProductRepository) Update(ctx context.Context, p domain.Product) (domain.Product, error) {
	updated, err := r.q.UpdateProduct(ctx, toUpdateParams(p))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, domain.ErrProductNotFound
	}
	if err != nil {
		return domain.Product{}, fmt.Errorf("repository: update product: %w", err)
	}
	return toDomain(updated), nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	affected, err := r.q.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("repository: delete product: %w", err)
	}
	if affected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}
