package repository

import (
	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/products/internal/db"
	"github.com/Jeudry/adventist-stack/services/products/internal/domain"
)

func toDomain(p db.Product) domain.Product {
	return domain.Product{
		Base: entity.Base{
			ID:        p.ID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			DeletedAt: p.DeletedAt,
			CreatedBy: p.CreatedBy,
			UpdatedBy: p.UpdatedBy,
			DeletedBy: p.DeletedBy,
		},
		Name:        p.Name,
		Sku:         p.Sku,
		Description: p.Description,
		Brand:       p.Brand,
		ReleaseDate: p.ReleaseDate,
		Status:      statusFromDB(p.Status),
	}
}

func statusFromDB(s int16) domain.Status {
	switch s {
	case 2:
		return domain.StatusInactive
	case 3:
		return domain.StatusDiscontinued
	default:
		return domain.StatusActive
	}
}

func statusToDB(s domain.Status) int16 {
	switch s {
	case domain.StatusInactive:
		return 2
	case domain.StatusDiscontinued:
		return 3
	default:
		return 1
	}
}

func toCreateParams(p domain.Product) db.CreateProductParams {
	return db.CreateProductParams{
		Name:        p.Name,
		Sku:         p.Sku,
		Description: p.Description,
		Brand:       p.Brand,
		ReleaseDate: p.ReleaseDate,
		Status:      statusToDB(p.Status),
	}
}

func toUpdateParams(p domain.Product) db.UpdateProductParams {
	return db.UpdateProductParams{
		ID:          p.ID,
		Name:        p.Name,
		Sku:         p.Sku,
		Description: p.Description,
		Brand:       p.Brand,
		ReleaseDate: p.ReleaseDate,
		Status:      statusToDB(p.Status),
	}
}

func toListParams(q pagination.Query) db.ListProductsParams {
	return db.ListProductsParams{
		Search:    q.Search,
		RowLimit:  int32(q.Limit),
		RowOffset: int32(q.Offset),
	}
}
