package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/service"
)

type Handler struct {
	svc *service.PrayerService
}

func NewHandler(svc *service.PrayerService) *Handler {
	return &Handler{svc: svc}
}

type prayerIDInput struct {
	ID string `path:"id" format:"uuid" doc:"Identificador de la petición"`
}

type updatePrayerInput struct {
	ID   string `path:"id" format:"uuid"`
	Body PrayerRequest
}

type listPrayersInput struct {
	Page     int    `query:"page" default:"1" minimum:"1" doc:"Arranca en 1"`
	PageSize int    `query:"pageSize" default:"20" minimum:"1" maximum:"100"`
	Search   string `query:"search"`
}

type deletedOutput struct{}

func (h *Handler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "createPrayer",
		Method:        "POST",
		Path:          "/api/v1/prayers",
		Summary:       "Crear una petición de oración",
		Tags:          []string{"prayers"},
		DefaultStatus: 201,
	}, h.Create)

	huma.Register(api, huma.Operation{
		OperationID: "listPrayers",
		Method:      "GET",
		Path:        "/api/v1/prayers",
		Summary:     "Listar peticiones de oración",
		Tags:        []string{"prayers"},
	}, h.List)

	huma.Register(api, huma.Operation{
		OperationID: "getPrayer",
		Method:      "GET",
		Path:        "/api/v1/prayers/{id}",
		Summary:     "Obtener una petición",
		Tags:        []string{"prayers"},
	}, h.GetByID)

	huma.Register(api, huma.Operation{
		OperationID: "updatePrayer",
		Method:      "PUT",
		Path:        "/api/v1/prayers/{id}",
		Summary:     "Actualizar una petición",
		Description: "Reemplaza el recurso completo; los campos ausentes se vacían.",
		Tags:        []string{"prayers"},
	}, h.Update)

	huma.Register(api, huma.Operation{
		OperationID:   "deletePrayer",
		Method:        "DELETE",
		Path:          "/api/v1/prayers/{id}",
		Summary:       "Eliminar una petición",
		Tags:          []string{"prayers"},
		DefaultStatus: 204,
	}, h.Delete)
}

func (h *Handler) Create(ctx context.Context, in *httpx.In[PrayerRequest]) (*httpx.Out[PrayerVM], error) {
	actor, err := httpx.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	created, err := h.svc.Create(ctx, toDomain(in.Body, uuid.Nil, actor))
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[PrayerVM]{Body: toPrayerVM(created)}, nil
}

func (h *Handler) GetByID(ctx context.Context, in *prayerIDInput) (*httpx.Out[PrayerVM], error) {
	found, err := h.svc.GetByID(ctx, uuid.MustParse(in.ID))
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[PrayerVM]{Body: toPrayerVM(found)}, nil
}

func (h *Handler) List(ctx context.Context, in *listPrayersInput) (*httpx.Out[httpx.PageResponse[PrayerVM]], error) {
	page, err := h.svc.RetrieveList(ctx, pagination.ListRequest{
		Page:     in.Page,
		PageSize: in.PageSize,
		Search:   in.Search,
	}.ToQuery())
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[httpx.PageResponse[PrayerVM]]{
		Body: httpx.ToPageResponse(page, toPrayerVM),
	}, nil
}

func (h *Handler) Update(ctx context.Context, in *updatePrayerInput) (*httpx.Out[PrayerVM], error) {
	actor, err := httpx.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	updated, err := h.svc.Update(ctx, toDomain(in.Body, uuid.MustParse(in.ID), actor))
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[PrayerVM]{Body: toPrayerVM(updated)}, nil
}

func (h *Handler) Delete(ctx context.Context, in *prayerIDInput) (*deletedOutput, error) {
	actor, err := httpx.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.svc.Delete(ctx, uuid.MustParse(in.ID), actor); err != nil {
		return nil, domainError(err)
	}
	return &deletedOutput{}, nil
}
