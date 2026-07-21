package repository

import (
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/products/internal/db"
	"github.com/Jeudry/adventist-stack/services/products/internal/domain"
)

func toDomain(p db.Product) domain.Product {
	return domain.Product{
		Id:          p.ID,
		Name:        p.Name,
		Sku:         p.Sku,
		Description: p.Description,
		Brand:       p.Brand,
		ReleaseDate: p.ReleaseDate,
		Status:      domain.Status(p.Status),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func toCreateParams(p domain.Product) db.CreateProductParams {
	return db.CreateProductParams{
		Name:        p.Name,
		Sku:         p.Sku,
		Description: p.Description,
		Brand:       p.Brand,
		ReleaseDate: p.ReleaseDate,
		Status:      string(p.Status),
	}
}

func toUpdateParams(p domain.Product) db.UpdateProductParams {
	return db.UpdateProductParams{
		ID:          p.Id,
		Name:        p.Name,
		Sku:         p.Sku,
		Description: p.Description,
		Brand:       p.Brand,
		ReleaseDate: p.ReleaseDate,
		Status:      string(p.Status),
	}
}

func toListParams(q pagination.Query) db.ListProductsParams {
	return db.ListProductsParams{
		Search:    q.Search,
		RowLimit:  int32(q.Limit),
		RowOffset: int32(q.Offset),
	}
}
