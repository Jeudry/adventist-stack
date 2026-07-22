package grpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	productsv1 "github.com/Jeudry/adventist-stack/gen/products/v1"
	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/products/internal/domain"
	"github.com/Jeudry/adventist-stack/services/products/internal/service"
	"github.com/google/uuid"
)

type Server struct {
	productsv1.UnimplementedProductServiceServer
	svc *service.ProductService
}

func NewServer(svc *service.ProductService) *Server {
	return &Server{svc: svc}
}

func (s *Server) CreateProduct(ctx context.Context, req *productsv1.CreateProductRequest) (*productsv1.CreateProductResponse, error) {
	created, err := s.svc.Create(ctx, productFromCreate(req))
	if err != nil {
		return nil, toStatus(err)
	}
	return &productsv1.CreateProductResponse{Product: productToProto(created)}, nil
}

func (s *Server) GetProduct(ctx context.Context, req *productsv1.GetProductRequest) (*productsv1.GetProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	found, err := s.svc.Get(ctx, id)
	if err != nil {
		return nil, toStatus(err)
	}
	return &productsv1.GetProductResponse{Product: productToProto(found)}, nil
}

func (s *Server) ListProducts(ctx context.Context, req *productsv1.ListProductsRequest) (*productsv1.ListProductsResponse, error) {
	page, err := s.svc.List(ctx, pagination.ListRequest{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
		Search:   deref(req.Search),
	})
	if err != nil {
		return nil, toStatus(err)
	}
	items := make([]*productsv1.Product, len(page.Items))
	for i, p := range page.Items {
		items[i] = productToProto(p)
	}
	return &productsv1.ListProductsResponse{
		Items:    items,
		Total:    int32(page.Total),
		Page:     int32(page.Page),
		PageSize: int32(page.PageSize),
	}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, req *productsv1.UpdateProductRequest) (*productsv1.UpdateProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	updated, err := s.svc.Update(ctx, productFromUpdate(req, id))
	if err != nil {
		return nil, toStatus(err)
	}
	return &productsv1.UpdateProductResponse{Product: productToProto(updated)}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *productsv1.DeleteProductRequest) (*productsv1.DeleteProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	if err := s.svc.Delete(ctx, id); err != nil {
		return nil, toStatus(err)
	}
	return &productsv1.DeleteProductResponse{Deleted: true}, nil
}

func productFromCreate(req *productsv1.CreateProductRequest) domain.Product {
	return domain.Product{
		Name:        req.Name,
		Sku:         req.Sku,
		Description: req.Description,
		Brand:       req.Brand,
		ReleaseDate: tsToPtr(req.ReleaseDate),
		Status:      statusFromProto(req.Status),
	}
}

func productFromUpdate(req *productsv1.UpdateProductRequest, id uuid.UUID) domain.Product {
	return domain.Product{
		Base:        entity.Base{ID: id},
		Name:        req.Name,
		Sku:         req.Sku,
		Description: req.Description,
		Brand:       req.Brand,
		ReleaseDate: tsToPtr(req.ReleaseDate),
		Status:      statusFromProto(req.Status),
	}
}

func productToProto(p domain.Product) *productsv1.Product {
	return &productsv1.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Sku:         p.Sku,
		Description: p.Description,
		Brand:       p.Brand,
		ReleaseDate: ptrToTs(p.ReleaseDate),
		Status:      statusToProto(p.Status),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}

func statusFromProto(s productsv1.ProductStatus) domain.Status {
	switch s {
	case productsv1.ProductStatus_PRODUCT_STATUS_INACTIVE:
		return domain.StatusInactive
	case productsv1.ProductStatus_PRODUCT_STATUS_DISCONTINUED:
		return domain.StatusDiscontinued
	default:
		return domain.StatusActive
	}
}

func statusToProto(s domain.Status) productsv1.ProductStatus {
	switch s {
	case domain.StatusInactive:
		return productsv1.ProductStatus_PRODUCT_STATUS_INACTIVE
	case domain.StatusDiscontinued:
		return productsv1.ProductStatus_PRODUCT_STATUS_DISCONTINUED
	default:
		return productsv1.ProductStatus_PRODUCT_STATUS_ACTIVE
	}
}

func tsToPtr(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func ptrToTs(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toStatus(err error) error {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidProduct):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
