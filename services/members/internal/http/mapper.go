package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/ptr"
	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
)

func toDomain(req MemberRequest, id, actor uuid.UUID) (domain.Member, error) {
	email, err := vo.NewOptionalEmail(req.Email)
	if err != nil {
		return domain.Member{}, huma.Error422UnprocessableEntity(err.Error())
	}
	phone, err := vo.NewOptionalPhone(req.Phone)
	if err != nil {
		return domain.Member{}, huma.Error422UnprocessableEntity(err.Error())
	}
	birthDate, err := httpx.ParseDate(req.BirthDate)
	if err != nil {
		return domain.Member{}, huma.Error422UnprocessableEntity("birthDate must be YYYY-MM-DD")
	}
	baptismDate, err := httpx.ParseDate(req.BaptismDate)
	if err != nil {
		return domain.Member{}, huma.Error422UnprocessableEntity("baptismDate must be YYYY-MM-DD")
	}

	return domain.Member{
		Base:        entity.Base{ID: id, CreatedBy: actor, UpdatedBy: ptr.NonZero(actor)},
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       email,
		Phone:       phone,
		Gender:      domain.ParseGender(req.Gender),
		Address:     req.Address,
		BirthDate:   birthDate,
		BaptismDate: baptismDate,
		Status:      domain.ParseStatus(req.Status),
	}, nil
}

func toMemberVM(m domain.Member) MemberVM {
	return MemberVM{
		BaseVM:      httpx.ToBaseVM(m.Base),
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Email:       m.Email.Ptr(),
		Phone:       m.Phone.Ptr(),
		Gender:      m.Gender.String(),
		Address:     m.Address,
		BirthDate:   httpx.FormatDate(m.BirthDate),
		BaptismDate: httpx.FormatDate(m.BaptismDate),
		Status:      m.Status.String(),
	}
}

func domainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrMemberNotFound):
		return huma.Error404NotFound(err.Error())
	case errors.Is(err, domain.ErrorInvalidMember),
		errors.Is(err, vo.ErrInvalidEmail),
		errors.Is(err, vo.ErrInvalidPhone):
		return huma.Error422UnprocessableEntity(err.Error())
	default:
		return huma.Error500InternalServerError("internal error")
	}
}
