package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	membersv1 "github.com/Jeudry/adventist-stack/gen/members/v1"
	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/pkg/ptr"
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
	"github.com/Jeudry/adventist-stack/services/members/internal/service"
	"github.com/google/uuid"
)

type Server struct {
	membersv1.UnimplementedMemberServiceServer
	svc *service.MemberService
}

func NewServer(svc *service.MemberService) *Server {
	return &Server{svc: svc}
}

func (s *Server) CreateMember(ctx context.Context, req *membersv1.CreateMemberRequest) (*membersv1.Member, error) {
	member, err := memberFromCreate(req)
	if err != nil {
		return nil, toStatus(err)
	}
	created, err := s.svc.Create(ctx, member)
	if err != nil {
		return nil, toStatus(err)
	}
	return memberToProto(created), nil
}

func (s *Server) GetMember(ctx context.Context, req *membersv1.GetMemberRequest) (*membersv1.Member, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	found, err := s.svc.GetByID(ctx, domain.Member{Base: entity.Base{ID: id}})
	if err != nil {
		return nil, toStatus(err)
	}
	return memberToProto(found), nil
}

func (s *Server) ListMembers(ctx context.Context, req *membersv1.ListMembersRequest) (*membersv1.ListMembersResponse, error) {
	page, err := s.svc.RetrieveList(ctx, pagination.ListRequest{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
		Search:   ptr.Deref(req.Search),
	})
	if err != nil {
		return nil, toStatus(err)
	}
	items := make([]*membersv1.Member, len(page.Items))
	for i, m := range page.Items {
		items[i] = memberToProto(m)
	}
	return &membersv1.ListMembersResponse{
		Items:    items,
		Total:    int32(page.Total),
		Page:     int32(page.Page),
		PageSize: int32(page.PageSize),
	}, nil
}

func (s *Server) UpdateMember(ctx context.Context, req *membersv1.UpdateMemberRequest) (*membersv1.Member, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	member, err := memberFromUpdate(req, id)
	if err != nil {
		return nil, toStatus(err)
	}
	updated, err := s.svc.Update(ctx, member)
	if err != nil {
		return nil, toStatus(err)
	}
	return memberToProto(updated), nil
}

func (s *Server) DeleteMember(ctx context.Context, req *membersv1.DeleteMemberRequest) (*membersv1.DeleteMemberResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	if err := s.svc.Delete(ctx, id); err != nil {
		return nil, toStatus(err)
	}
	return &membersv1.DeleteMemberResponse{Deleted: true}, nil
}
