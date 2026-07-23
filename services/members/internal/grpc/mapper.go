package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	membersv1 "github.com/Jeudry/adventist-stack/gen/members/v1"
	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/protoconv"
	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
	"github.com/google/uuid"
)

func memberFromCreate(req *membersv1.CreateMemberRequest) (domain.Member, error) {
	email, err := vo.NewOptionalEmail(req.Email)
	if err != nil {
		return domain.Member{}, err
	}
	phone, err := vo.NewOptionalPhone(req.Phone)
	if err != nil {
		return domain.Member{}, err
	}

	return domain.Member{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       email,
		Phone:       phone,
		Gender:      domain.ParseGender(req.Gender),
		Address:     req.Address,
		BirthDate:   protoconv.TimeFromProto(req.BirthDate),
		BaptismDate: protoconv.TimeFromProto(req.BaptismDate),
		Status:      statusFromProto(req.Status),
	}, nil
}

func memberFromUpdate(req *membersv1.UpdateMemberRequest, id uuid.UUID) (domain.Member, error) {
	email, err := vo.NewOptionalEmail(req.Email)
	if err != nil {
		return domain.Member{}, err
	}
	phone, err := vo.NewOptionalPhone(req.Phone)
	if err != nil {
		return domain.Member{}, err
	}

	return domain.Member{
		Base:        entity.Base{ID: id},
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       email,
		Phone:       phone,
		Gender:      domain.ParseGender(req.Gender),
		Address:     req.Address,
		BirthDate:   protoconv.TimeFromProto(req.BirthDate),
		BaptismDate: protoconv.TimeFromProto(req.BaptismDate),
		Status:      statusFromProto(req.Status),
	}, nil
}

func memberToProto(m domain.Member) *membersv1.Member {
	return &membersv1.Member{
		Id:          m.ID.String(),
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Email:       m.Email.Ptr(),
		Phone:       m.Phone.Ptr(),
		Gender:      m.Gender.String(),
		Address:     m.Address,
		BirthDate:   protoconv.TimeToProto(m.BirthDate),
		BaptismDate: protoconv.TimeToProto(m.BaptismDate),
		Status:      statusToProto(m.Status),
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
	}
}

func statusFromProto(s membersv1.MemberStatus) domain.Status {
	switch s {
	case membersv1.MemberStatus_MEMBER_STATUS_INACTIVE:
		return domain.StatusInactive
	case membersv1.MemberStatus_MEMBER_STATUS_VISITOR:
		return domain.StatusVisitor
	default:
		return domain.StatusActive
	}
}

func statusToProto(s domain.Status) membersv1.MemberStatus {
	switch s {
	case domain.StatusInactive:
		return membersv1.MemberStatus_MEMBER_STATUS_INACTIVE
	case domain.StatusVisitor:
		return membersv1.MemberStatus_MEMBER_STATUS_VISITOR
	default:
		return membersv1.MemberStatus_MEMBER_STATUS_ACTIVE
	}
}

func toStatus(err error) error {
	switch {
	case errors.Is(err, domain.ErrMemberNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrorInvalidMember),
		errors.Is(err, vo.ErrInvalidEmail),
		errors.Is(err, vo.ErrInvalidPhone):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
