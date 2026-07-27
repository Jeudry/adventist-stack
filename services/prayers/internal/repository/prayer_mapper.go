package repository

import (
	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/db"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/domain"
)

func toDomain(m db.Prayer) (domain.Prayer, error) {
	return domain.Prayer{
		Title:       m.Title,
		Description: m.Description,
		AuthorName:  *m.AuthorName,
		IsAnonymous: m.IsAnonymous,
		Status:      domain.Status(m.Status),
		Base: entity.Base{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			CreatedBy: m.CreatedBy,
			UpdatedAt: m.UpdatedAt,
			UpdatedBy: m.UpdatedBy,
			DeletedAt: m.DeletedAt,
			DeletedBy: m.DeletedBy,
		},
	}, nil
}

func toCreateParams(p domain.Prayer) db.CreatePayerParams {
	var authorName *string
	if p.AuthorName != "" {
		authorName = &p.AuthorName
	}

	return db.CreatePayerParams{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		AuthorName:  authorName,
		IsAnonymous: p.IsAnonymous,
		Status:      int32(p.Status),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func toUpdateParams(p domain.Prayer) db.UpdatePrayerParams {
	var authorName *string
	if p.AuthorName != "" {
		authorName = &p.AuthorName
	}

	return db.UpdatePrayerParams{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		AuthorName:  authorName,
		IsAnonymous: p.IsAnonymous,
		Status:      int32(p.Status),
		UpdatedBy:   p.UpdatedBy,
	}
}

func toListParams(q pagination.Query) db.ListPrayersParams {
	return db.ListPrayersParams{
		Search:    q.Search,
		RowOffset: int32(q.Offset),
		RowLimit:  int32(q.Limit),
	}
}
