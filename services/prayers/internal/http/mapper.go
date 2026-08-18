package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/ptr"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/domain"
)

// isAnonymous stays in the JSON because that is how a client says what it wants;
// inside, it collapses into the one field that decides it. Asking for anonymity
// drops the name here, so it never reaches the database.
func toDomain(req PrayerRequest, id, actor uuid.UUID) domain.Prayer {
	var authorName *string
	if !req.IsAnonymous && req.AuthorName != "" {
		authorName = &req.AuthorName
	}

	return domain.Prayer{
		Base:        entity.Base{ID: id, CreatedBy: actor, UpdatedBy: ptr.NonZero(actor)},
		Title:       req.Title,
		Description: req.Description,
		AuthorName:  authorName,
		Status:      domain.ParseStatus(req.Status),
	}
}

func toPrayerVM(p domain.Prayer) PrayerVM {
	return PrayerVM{
		BaseVM:      httpx.ToBaseVM(p.Base),
		Title:       p.Title,
		Description: p.Description,
		AuthorName:  ptr.Deref(p.AuthorName),
		IsAnonymous: p.AuthorName == nil,
		Status:      p.Status.String(),
	}
}

func domainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrPrayerNotFound):
		return huma.Error404NotFound(err.Error())
	case errors.Is(err, domain.ErrInvalidPrayer), errors.Is(err, domain.ErrInvalidStatus):
		return huma.Error422UnprocessableEntity(err.Error())
	default:
		return huma.Error500InternalServerError("internal error")
	}
}
