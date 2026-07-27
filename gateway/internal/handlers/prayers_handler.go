package handlers

import (
	"net/http"

	"github.com/Jeudry/adventist-stack/gateway/internal/mappers"
	"github.com/Jeudry/adventist-stack/gateway/internal/models/base"
	"github.com/Jeudry/adventist-stack/gateway/internal/models/prayer"
	prayersv1 "github.com/Jeudry/adventist-stack/gen/prayers/v1"
)

type PrayersHandler struct {
	client prayersv1.PrayerServiceClient
}

func NewPrayersHandler(client prayersv1.PrayerServiceClient) *PrayersHandler {
	return &PrayersHandler{
		client: client,
	}
}

func (h *PrayersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req prayer.CreatePrayerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}

	protoReq := mappers.ToCreatePrayerProto(req)

	res, err := h.client.CreatePrayer(r.Context(), protoReq)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *PrayersHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	protoReq := &prayersv1.GetPrayerRequest{Id: id}
	res, err := h.client.GetPrayer(r.Context(), protoReq)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *PrayersHandler) List(w http.ResponseWriter, r *http.Request) {
	query := parseListQuery(r)
	protoReq := &prayersv1.ListPrayersRequest{
		Page:     query.Page,
		PageSize: query.PageSize,
		Search:   &query.Search,
	}

	res, err := h.client.ListPrayers(r.Context(), protoReq)
	if err != nil {
		writeError(w, err)
		return
	}

	items := make([]prayer.PrayerVM, len(res.GetItems()))
	for i, item := range res.GetItems() {
		items[i] = mappers.ToPrayerVM(item)
	}

	pageResp := base.PageResponse[prayer.PrayerVM]{
		Items:    items,
		Total:    res.GetTotal(),
		Page:     res.GetPage(),
		PageSize: res.GetPageSize(),
	}

	writeJSON(w, http.StatusOK, pageResp)
}

func (h *PrayersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	var req prayer.UpdatePrayerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}

	protoReq := mappers.ToUpdatePrayerProto(id, req)
	res, err := h.client.UpdatePrayer(r.Context(), protoReq)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *PrayersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	protoReq := &prayersv1.DeletePrayerRequest{Id: id}
	res, err := h.client.DeletePrayer(r.Context(), protoReq)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}
