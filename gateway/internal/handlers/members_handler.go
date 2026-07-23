package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	membersv1 "github.com/Jeudry/adventist-stack/gen/members/v1"
	"github.com/go-chi/chi/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MembersHandler struct {
	client membersv1.MemberServiceClient
}

func NewMembersHandler(client membersv1.MemberServiceClient) *MembersHandler {
	return &MembersHandler{
		client: client,
	}
}

type CreateMemberRequest struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Gender      string  `json:"gender"`
	Address     *string `json:"address,omitempty"`
	BirthDate   *string `json:"birth_date,omitempty"`   // Formato "YYYY-MM-DD"
	BaptismDate *string `json:"baptism_date,omitempty"` // Formato "YYYY-MM-DD"
	Status      string  `json:"status,omitempty"`
}

type UpdateMemberRequest struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Gender      string  `json:"gender"`
	Address     *string `json:"address,omitempty"`
	BirthDate   *string `json:"birth_date,omitempty"`   // Formato "YYYY-MM-DD"
	BaptismDate *string `json:"baptism_date,omitempty"` // Formato "YYYY-MM-DD"
	Status      string  `json:"status,omitempty"`
}

type MemberVM struct {
	BaseVM
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Gender      string  `json:"gender"`
	Address     *string `json:"address,omitempty"`
	BirthDate   *string `json:"birth_date,omitempty"`   // Formato "YYYY-MM-DD"
	BaptismDate *string `json:"baptism_date,omitempty"` // Formato "YYYY-MM-DD"
	Status      string  `json:"status"`
}

func (h *MembersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateMemberRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}

	birthDate, err := parseDatePtr(req.BirthDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid birth_date format, use YYYY-MM-DD"})
		return
	}
	baptismDate, err := parseDatePtr(req.BaptismDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid baptism_date format, use YYYY-MM-DD"})
		return
	}

	res, err := h.client.CreateMember(r.Context(), &membersv1.CreateMemberRequest{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		Gender:      req.Gender,
		Address:     req.Address,
		BirthDate:   birthDate,
		BaptismDate: baptismDate,
		Status:      statusToProto(req.Status),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toMemberVM(res.GetMember()))
}

func (h *MembersHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := h.client.GetMember(r.Context(), &membersv1.GetMemberRequest{Id: id})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toMemberVM(res.GetMember()))
}

func (h *MembersHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	search := r.URL.Query().Get("search")

	res, err := h.client.ListMembers(r.Context(), &membersv1.ListMembersRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   &search,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	items := make([]MemberVM, len(res.GetItems()))
	for i, m := range res.GetItems() {
		items[i] = toMemberVM(m)
	}

	writeJSON(w, http.StatusOK, PageResponse[MemberVM]{
		Items:    items,
		Total:    res.GetTotal(),
		Page:     res.GetPage(),
		PageSize: res.GetPageSize(),
	})
}

func (h *MembersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateMemberRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}

	birthDate, err := parseDatePtr(req.BirthDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid birth_date format, use YYYY-MM-DD"})
		return
	}
	baptismDate, err := parseDatePtr(req.BaptismDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid baptism_date format, use YYYY-MM-DD"})
		return
	}

	res, err := h.client.UpdateMember(r.Context(), &membersv1.UpdateMemberRequest{
		Id:          id,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		Gender:      req.Gender,
		Address:     req.Address,
		BirthDate:   birthDate,
		BaptismDate: baptismDate,
		Status:      statusToProto(req.Status),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toMemberVM(res.GetMember()))
}

func (h *MembersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.client.DeleteMember(r.Context(), &membersv1.DeleteMemberRequest{Id: id})
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toMemberVM(m *membersv1.Member) MemberVM {
	if m == nil {
		return MemberVM{}
	}

	return MemberVM{
		BaseVM:      toBaseVM(m.GetId(), m.GetCreatedAt(), m.GetUpdatedAt()),
		FirstName:   m.GetFirstName(),
		LastName:    m.GetLastName(),
		Email:       m.Email,
		Phone:       m.Phone,
		Gender:      m.GetGender(),
		Address:     m.Address,
		BirthDate:   formatDatePtr(m.GetBirthDate()),
		BaptismDate: formatDatePtr(m.GetBaptismDate()),
		Status:      m.GetStatus().String(),
	}
}

func parseDatePtr(s *string) (*timestamppb.Timestamp, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, err
	}
	return timestamppb.New(t), nil
}

func formatDatePtr(ts *timestamppb.Timestamp) *string {
	if ts == nil {
		return nil
	}
	s := ts.AsTime().Format("2006-01-02")
	return &s
}

func statusToProto(s string) membersv1.MemberStatus {
	switch strings.ToLower(s) {
	case "inactive":
		return membersv1.MemberStatus_MEMBER_STATUS_INACTIVE
	case "visitor":
		return membersv1.MemberStatus_MEMBER_STATUS_VISITOR
	default:
		return membersv1.MemberStatus_MEMBER_STATUS_ACTIVE
	}
}
