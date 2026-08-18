package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/ptr"
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/domain"
)

// actor is the caller the gateway authenticated. Both audit fields are filled on every crossing
// because each statement reads only its own: the insert takes created_by and never sees
// updated_by, the update takes updated_by and never touches created_by.
func toDomain(req SabbathSchoolRequest, id, actor uuid.UUID) (domain.SabbathSchool, error) {
	teacherID, err := parseOptionalUUID(req.TeacherID)
	if err != nil {
		return domain.SabbathSchool{}, huma.Error422UnprocessableEntity("teacherId is not a uuid")
	}

	return domain.SabbathSchool{
		Base:         entity.Base{ID: id, CreatedBy: actor, UpdatedBy: ptr.NonZero(actor)},
		Name:         req.Name,
		TeacherID:    teacherID,
		Location:     req.Location,
		TargetMinAge: req.TargetMinAge,
		TargetMaxAge: req.TargetMaxAge,
		Status:       domain.ParseStatus(req.Status),
	}, nil
}

func parseOptionalUUID(raw *string) (uuid.NullUUID, error) {
	if raw == nil || *raw == "" {
		return uuid.NullUUID{}, nil
	}
	parsed, err := uuid.Parse(*raw)
	if err != nil {
		return uuid.NullUUID{}, err
	}
	return uuid.NullUUID{UUID: parsed, Valid: true}, nil
}

func toSabbathSchoolVM(ss domain.SabbathSchool) SabbathSchoolVM {
	var teacherID *string
	if ss.TeacherID.Valid {
		id := ss.TeacherID.UUID.String()
		teacherID = &id
	}
	return SabbathSchoolVM{
		BaseVM:       httpx.ToBaseVM(ss.Base),
		Name:         ss.Name,
		TeacherID:    teacherID,
		Location:     ss.Location,
		TargetMinAge: ss.TargetMinAge,
		TargetMaxAge: ss.TargetMaxAge,
		Status:       ss.Status.String(),
	}
}

// domainError maps what the domain refuses onto huma's problem responses, which
// carry the status inside the body as RFC 7807 rather than a bare {"error"}.
func domainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrSabbathSchoolNotFound):
		return huma.Error404NotFound(err.Error())
	case errors.Is(err, domain.ErrInvalidSabbathSchool), errors.Is(err, domain.ErrInvalidStatus):
		return huma.Error422UnprocessableEntity(err.Error())
	default:
		return huma.Error500InternalServerError("internal error")
	}
}
