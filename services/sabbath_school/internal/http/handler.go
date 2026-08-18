package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/service"
)

type Handler struct {
	svc *service.SabbathSchoolService
}

func NewHandler(svc *service.SabbathSchoolService) *Handler {
	return &Handler{svc: svc}
}

type idInput struct {
	ID string `path:"id" format:"uuid" doc:"Identificador de la escuela sabática"`
}

type updateSabbathSchoolInput struct {
	ID   string `path:"id" format:"uuid" doc:"Identificador de la escuela sabática"`
	Body SabbathSchoolRequest
}

type listSabbathSchoolsInput struct {
	Page     int    `query:"page" default:"1" minimum:"1" doc:"Arranca en 1"`
	PageSize int    `query:"pageSize" default:"20" minimum:"1" maximum:"100"`
	Search   string `query:"search"`
}

// Register declares the operations. What used to live in the sabbath-schools
// block of openapi.yaml lives here now, and the spec is generated from it.
func (h *Handler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "createSabbathSchool",
		Method:        "POST",
		Path:          "/api/v1/sabbath-schools",
		Summary:       "Crear una escuela sabática",
		Tags:          []string{"sabbath-schools"},
		DefaultStatus: 201,
	}, h.Create)

	huma.Register(api, huma.Operation{
		OperationID: "listSabbathSchools",
		Method:      "GET",
		Path:        "/api/v1/sabbath-schools",
		Summary:     "Listar escuelas sabáticas",
		Tags:        []string{"sabbath-schools"},
	}, h.List)

	huma.Register(api, huma.Operation{
		OperationID: "getSabbathSchool",
		Method:      "GET",
		Path:        "/api/v1/sabbath-schools/{id}",
		Summary:     "Obtener una escuela sabática",
		Tags:        []string{"sabbath-schools"},
	}, h.GetByID)

	huma.Register(api, huma.Operation{
		OperationID: "updateSabbathSchool",
		Method:      "PUT",
		Path:        "/api/v1/sabbath-schools/{id}",
		Summary:     "Actualizar una escuela sabática",
		Tags:        []string{"sabbath-schools"},
	}, h.Update)

	huma.Register(api, huma.Operation{
		OperationID:   "deleteSabbathSchool",
		Method:        "DELETE",
		Path:          "/api/v1/sabbath-schools/{id}",
		Summary:       "Eliminar una escuela sabática",
		Tags:          []string{"sabbath-schools"},
		DefaultStatus: 204,
	}, h.Delete)
}

func (h *Handler) Create(ctx context.Context, in *httpx.In[SabbathSchoolRequest]) (*httpx.Out[SabbathSchoolVM], error) {
	sabbathSchool, err := toDomain(in.Body, uuid.Nil, httpx.UserID(ctx))
	if err != nil {
		return nil, err
	}

	created, err := h.svc.Create(ctx, sabbathSchool)
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[SabbathSchoolVM]{Body: toSabbathSchoolVM(created)}, nil
}

func (h *Handler) GetByID(ctx context.Context, in *idInput) (*httpx.Out[SabbathSchoolVM], error) {
	// The uuid format tag already rejected anything malformed, so this parse
	// cannot fail by the time the handler runs.
	found, err := h.svc.GetByID(ctx, uuid.MustParse(in.ID))
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[SabbathSchoolVM]{Body: toSabbathSchoolVM(found)}, nil
}

func (h *Handler) List(ctx context.Context, in *listSabbathSchoolsInput) (*httpx.Out[httpx.PageResponse[SabbathSchoolVM]], error) {
	page, err := h.svc.RetrieveList(ctx, pagination.ListRequest{
		Page:     in.Page,
		PageSize: in.PageSize,
		Search:   in.Search,
	}.ToQuery())
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[httpx.PageResponse[SabbathSchoolVM]]{
		Body: httpx.ToPageResponse(page, toSabbathSchoolVM),
	}, nil
}

func (h *Handler) Update(ctx context.Context, in *updateSabbathSchoolInput) (*httpx.Out[SabbathSchoolVM], error) {
	sc, err := toDomain(in.Body, uuid.MustParse(in.ID), httpx.UserID(ctx))
	if err != nil {
		return nil, err
	}
	updated, err := h.svc.Update(ctx, sc)
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[SabbathSchoolVM]{Body: toSabbathSchoolVM(updated)}, nil
}

func (h *Handler) Delete(ctx context.Context, in *idInput) (*httpx.Out[struct{}], error) {
	if err := h.svc.Delete(ctx, uuid.MustParse(in.ID), httpx.UserID(ctx)); err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[struct{}]{}, nil
}
