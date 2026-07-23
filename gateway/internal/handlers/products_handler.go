package handlers

import (
	"net/http"
	"strconv"

	productsv1 "github.com/Jeudry/adventist-stack/gen/products/v1"
	"github.com/go-chi/chi/v5"
)

type ProductsHandler struct {
	client productsv1.ProductServiceClient
}

func NewProductsHandler(client productsv1.ProductServiceClient) *ProductsHandler {
	return &ProductsHandler{client: client}
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Sku         string  `json:"sku"`
	Description *string `json:"description,omitempty"`
	Brand       *string `json:"brand,omitempty"`
	Status      string  `json:"status,omitempty"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Sku         string  `json:"sku"`
	Description *string `json:"description,omitempty"`
	Brand       *string `json:"brand,omitempty"`
	Status      string  `json:"status,omitempty"`
}

type ProductVM struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Sku         string  `json:"sku"`
	Description *string `json:"description,omitempty"`
	Brand       *string `json:"brand,omitempty"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ProductListResponse struct {
	Items    []ProductVM `json:"items"`
	Total    int32       `json:"total"`
	Page     int32       `json:"page"`
	PageSize int32       `json:"page_size"`
}

func (h *ProductsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}

	res, err := h.client.CreateProduct(r.Context(), &productsv1.CreateProductRequest{
		Name:        req.Name,
		Sku:         req.Sku,
		Description: req.Description,
		Brand:       req.Brand,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toProductVM(res.GetProduct()))
}

func (h *ProductsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := h.client.GetProduct(r.Context(), &productsv1.GetProductRequest{Id: id})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toProductVM(res.GetProduct()))
}

func (h *ProductsHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	search := r.URL.Query().Get("search")

	res, err := h.client.ListProducts(r.Context(), &productsv1.ListProductsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   &search,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	items := make([]ProductVM, len(res.GetItems()))
	for i, p := range res.GetItems() {
		items[i] = toProductVM(p)
	}

	writeJSON(w, http.StatusOK, ProductListResponse{
		Items:    items,
		Total:    res.GetTotal(),
		Page:     res.GetPage(),
		PageSize: res.GetPageSize(),
	})
}

func (h *ProductsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateProductRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}

	res, err := h.client.UpdateProduct(r.Context(), &productsv1.UpdateProductRequest{
		Id:          id,
		Name:        req.Name,
		Sku:         req.Sku,
		Description: req.Description,
		Brand:       req.Brand,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toProductVM(res.GetProduct()))
}

func (h *ProductsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.client.DeleteProduct(r.Context(), &productsv1.DeleteProductRequest{Id: id})
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toProductVM(p *productsv1.Product) ProductVM {
	if p == nil {
		return ProductVM{}
	}
	var createdAt, updatedAt string
	if p.GetCreatedAt() != nil {
		createdAt = p.GetCreatedAt().AsTime().Format("2006-01-02T15:04:05Z07:00")
	}
	if p.GetUpdatedAt() != nil {
		updatedAt = p.GetUpdatedAt().AsTime().Format("2006-01-02T15:04:05Z07:00")
	}

	return ProductVM{
		ID:          p.GetId(),
		Name:        p.GetName(),
		Sku:         p.GetSku(),
		Description: p.Description,
		Brand:       p.Brand,
		Status:      p.GetStatus().String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
