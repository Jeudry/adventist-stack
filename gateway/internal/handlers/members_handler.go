package handlers

import (
	"net/http"

	membersv1 "github.com/Jeudry/adventist-stack/gen/members/v1"
	"github.com/Jeudry/adventist-stack/gateway/internal/mappers"
	"github.com/Jeudry/adventist-stack/gateway/internal/models/base"
	"github.com/Jeudry/adventist-stack/gateway/internal/models/member"
)

type MembersHandler struct {
	client membersv1.MemberServiceClient
}

func NewMembersHandler(client membersv1.MemberServiceClient) *MembersHandler {
	return &MembersHandler{
		client: client,
	}
}

func (h *MembersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req member.CreateMemberRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}

	protoReq, err := mappers.ToCreateMemberProto(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	res, err := h.client.CreateMember(r.Context(), protoReq)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, mappers.ToMemberVM(res))
}

func (h *MembersHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	protoReq := &membersv1.GetMemberRequest{Id: id}
	res, err := h.client.GetMember(r.Context(), protoReq)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mappers.ToMemberVM(res))
}

func (h *MembersHandler) List(w http.ResponseWriter, r *http.Request) {
	query := parseListQuery(r)

	protoReq := &membersv1.ListMembersRequest{
		Page:     query.Page,
		PageSize: query.PageSize,
		Search:   &query.Search,
	}

	res, err := h.client.ListMembers(r.Context(), protoReq)
	if err != nil {
		writeError(w, err)
		return
	}

	items := make([]member.MemberVM, len(res.GetItems()))
	for i, m := range res.GetItems() {
		items[i] = mappers.ToMemberVM(m)
	}

	pageResp := base.PageResponse[member.MemberVM]{
		Items:    items,
		Total:    res.GetTotal(),
		Page:     res.GetPage(),
		PageSize: res.GetPageSize(),
	}

	writeJSON(w, http.StatusOK, pageResp)
}

func (h *MembersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	var req member.UpdateMemberRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}

	protoReq, err := mappers.ToUpdateMemberProto(id, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	res, err := h.client.UpdateMember(r.Context(), protoReq)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mappers.ToMemberVM(res))
}

func (h *MembersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	protoReq := &membersv1.DeleteMemberRequest{Id: id}
	_, err := h.client.DeleteMember(r.Context(), protoReq)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
