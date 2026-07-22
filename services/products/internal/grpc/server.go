package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	productsv1 "github.com/Jeudry/adventist-stack/gen/products/v1"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/pkg/ptr"
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
		Search:   ptr.Deref(req.Search),
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
