package grpc

import (
	"context"

	prayersv1 "github.com/Jeudry/adventist-stack/gen/prayers/v1"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/pkg/ptr"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	prayersv1.UnimplementedPrayerServiceServer
	svc *service.PrayerService
}

func NewServer(svc *service.PrayerService) *Server {
	return &Server{svc: svc}
}

func (s *Server) CreatePrayer(ctx context.Context, req *prayersv1.CreatePrayerRequest) (*prayersv1.Prayer, error) {
	prayer, err := prayerFromCreateRequest(req)
	if err != nil {
		return nil, toStatus(err)
	}
	_, err = s.svc.Create(ctx, prayer)
	if err != nil {
		return nil, toStatus(err)
	}
	return prayerToProto(prayer), nil
}

func (s *Server) GetPrayer(ctx context.Context, req *prayersv1.GetPrayerRequest) (*prayersv1.Prayer, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	prayer, err := s.svc.GetByID(ctx, id)
	if err != nil {
		return nil, toStatus(err)
	}
	return prayerToProto(prayer), nil
}

func (s *Server) ListPrayers(ctx context.Context, req *prayersv1.ListPrayersRequest) (*prayersv1.ListPrayersResponse, error) {
	paginationReq := pagination.ListRequest{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
		Search:   ptr.Deref(req.Search),
	}.ToQuery()
	page, err := s.svc.RetrieveList(ctx, paginationReq)
	if err != nil {
		return nil, toStatus(err)
	}
	items := make([]*prayersv1.Prayer, len(page.Items))
	for _, item := range page.Items {
		items = append(items, prayerToProto(item))
	}
	return &prayersv1.ListPrayersResponse{
		Items:    items,
		Total:    int32(page.Total),
		Page:     int32(page.Page),
		PageSize: int32(page.PageSize),
	}, nil
}

func (s *Server) UpdatePrayer(ctx context.Context, req *prayersv1.UpdatePrayerRequest) (*prayersv1.Prayer, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	prayerUpdate, err := prayerFromUpdateRequest(req)
	if err != nil {
		return nil, toStatus(err)
	}
	prayerUpdate.ID = id
	prayer, err := s.svc.Update(ctx, prayerUpdate)
	if err != nil {
		return nil, toStatus(err)
	}
	return prayerToProto(prayer), nil
}

func (s *Server) DeletePrayer(ctx context.Context, req *prayersv1.DeletePrayerRequest) (*prayersv1.DeletePrayerResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	if err := s.svc.Delete(ctx, id); err != nil {
		return nil, toStatus(err)
	}
	return &prayersv1.DeletePrayerResponse{}, nil
}
