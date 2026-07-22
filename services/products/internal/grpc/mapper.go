package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	productsv1 "github.com/Jeudry/adventist-stack/gen/products/v1"
	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/protoconv"
	"github.com/Jeudry/adventist-stack/services/products/internal/domain"
	"github.com/google/uuid"
)

func productFromCreate(req *productsv1.CreateProductRequest) domain.Product {
	return domain.Product{
		Name:        req.Name,
		Sku:         req.Sku,
		Description: req.Description,
		Brand:       req.Brand,
		ReleaseDate: protoconv.TimeFromProto(req.ReleaseDate),
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
		ReleaseDate: protoconv.TimeFromProto(req.ReleaseDate),
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
		ReleaseDate: protoconv.TimeToProto(p.ReleaseDate),
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
