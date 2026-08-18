package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
	"github.com/Jeudry/adventist-stack/services/members/internal/service"
)

type Handler struct {
	svc *service.MemberService
}

func NewHandler(svc *service.MemberService) *Handler {
	return &Handler{svc: svc}
}

type memberIDInput struct {
	ID string `path:"id" format:"uuid" doc:"Identificador del miembro"`
}

type updateMemberInput struct {
	ID   string `path:"id" format:"uuid"`
	Body MemberRequest
}

type listMembersInput struct {
	Page     int    `query:"page" default:"1" minimum:"1" doc:"Arranca en 1"`
	PageSize int    `query:"pageSize" default:"20" minimum:"1" maximum:"100"`
	Search   string `query:"search"`
}

type deletedOutput struct{}

func (h *Handler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "createMember",
		Method:        "POST",
		Path:          "/api/v1/members",
		Summary:       "Crear un miembro",
		Tags:          []string{"members"},
		DefaultStatus: 201,
	}, h.Create)

	huma.Register(api, huma.Operation{
		OperationID: "listMembers",
		Method:      "GET",
		Path:        "/api/v1/members",
		Summary:     "Listar miembros",
		Tags:        []string{"members"},
	}, h.List)

	huma.Register(api, huma.Operation{
		OperationID: "getMember",
		Method:      "GET",
		Path:        "/api/v1/members/{id}",
		Summary:     "Obtener un miembro",
		Tags:        []string{"members"},
	}, h.GetByID)

	huma.Register(api, huma.Operation{
		OperationID: "updateMember",
		Method:      "PUT",
		Path:        "/api/v1/members/{id}",
		Summary:     "Actualizar un miembro",
		Description: "Reemplaza el recurso completo; los campos ausentes se vacían.",
		Tags:        []string{"members"},
	}, h.Update)

	huma.Register(api, huma.Operation{
		OperationID:   "deleteMember",
		Method:        "DELETE",
		Path:          "/api/v1/members/{id}",
		Summary:       "Eliminar un miembro",
		Tags:          []string{"members"},
		DefaultStatus: 204,
	}, h.Delete)
}

func (h *Handler) Create(ctx context.Context, in *httpx.In[MemberRequest]) (*httpx.Out[MemberVM], error) {
	member, err := toDomain(in.Body, uuid.Nil, httpx.UserID(ctx))
	if err != nil {
		return nil, err
	}

	created, err := h.svc.Create(ctx, member)
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[MemberVM]{Body: toMemberVM(created)}, nil
}

func (h *Handler) GetByID(ctx context.Context, in *memberIDInput) (*httpx.Out[MemberVM], error) {
	found, err := h.svc.GetByID(ctx, domain.Member{Base: entity.Base{ID: uuid.MustParse(in.ID)}})
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[MemberVM]{Body: toMemberVM(found)}, nil
}

func (h *Handler) List(ctx context.Context, in *listMembersInput) (*httpx.Out[httpx.PageResponse[MemberVM]], error) {
	page, err := h.svc.RetrieveList(ctx, pagination.ListRequest{
		Page:     in.Page,
		PageSize: in.PageSize,
		Search:   in.Search,
	}.ToQuery())
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[httpx.PageResponse[MemberVM]]{
		Body: httpx.ToPageResponse(page, toMemberVM),
	}, nil
}

func (h *Handler) Update(ctx context.Context, in *updateMemberInput) (*httpx.Out[MemberVM], error) {
	member, err := toDomain(in.Body, uuid.MustParse(in.ID), httpx.UserID(ctx))
	if err != nil {
		return nil, err
	}

	updated, err := h.svc.Update(ctx, member)
	if err != nil {
		return nil, domainError(err)
	}
	return &httpx.Out[MemberVM]{Body: toMemberVM(updated)}, nil
}

func (h *Handler) Delete(ctx context.Context, in *memberIDInput) (*deletedOutput, error) {
	if err := h.svc.Delete(ctx, uuid.MustParse(in.ID), httpx.UserID(ctx)); err != nil {
		return nil, domainError(err)
	}
	return &deletedOutput{}, nil
}
