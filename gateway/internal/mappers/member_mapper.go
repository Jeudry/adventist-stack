package mappers

import (
	"strings"
	"time"

	membersv1 "github.com/Jeudry/adventist-stack/gen/members/v1"
	"github.com/Jeudry/adventist-stack/gateway/internal/models/base"
	"github.com/Jeudry/adventist-stack/gateway/internal/models/member"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToMemberVM(m *membersv1.Member) member.MemberVM {
	if m == nil {
		return member.MemberVM{}
	}

	return member.MemberVM{
		BaseVM:      base.ToBaseVM(m.GetId(), m.GetCreatedAt(), m.GetUpdatedAt()),
		FirstName:   m.GetFirstName(),
		LastName:    m.LastName,
		Email:       m.Email,
		Phone:       m.Phone,
		Gender:      m.GetGender(),
		Address:     m.Address,
		BirthDate:   FormatDatePtr(m.GetBirthDate()),
		BaptismDate: FormatDatePtr(m.GetBaptismDate()),
		Status:      m.GetStatus().String(),
	}
}

func ParseDatePtr(s *string) (*timestamppb.Timestamp, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, err
	}
	return timestamppb.New(t), nil
}

func FormatDatePtr(ts *timestamppb.Timestamp) *string {
	if ts == nil {
		return nil
	}
	s := ts.AsTime().Format("2006-01-02")
	return &s
}

func ToCreateMemberProto(req member.CreateMemberRequest) (*membersv1.CreateMemberRequest, error) {
	birthDate, err := ParseDatePtr(req.BirthDate)
	if err != nil {
		return nil, err
	}
	baptismDate, err := ParseDatePtr(req.BaptismDate)
	if err != nil {
		return nil, err
	}

	return &membersv1.CreateMemberRequest{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		Gender:      req.Gender,
		Address:     req.Address,
		BirthDate:   birthDate,
		BaptismDate: baptismDate,
		Status:      MemberStatusToProto(req.Status),
	}, nil
}

func ToUpdateMemberProto(id string, req member.UpdateMemberRequest) (*membersv1.UpdateMemberRequest, error) {
	birthDate, err := ParseDatePtr(req.BirthDate)
	if err != nil {
		return nil, err
	}
	baptismDate, err := ParseDatePtr(req.BaptismDate)
	if err != nil {
		return nil, err
	}

	return &membersv1.UpdateMemberRequest{
		Id:          id,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		Gender:      req.Gender,
		Address:     req.Address,
		BirthDate:   birthDate,
		BaptismDate: baptismDate,
		Status:      MemberStatusToProto(req.Status),
	}, nil
}

func MemberStatusToProto(s string) membersv1.MemberStatus {
	switch strings.ToLower(s) {
	case "inactive":
		return membersv1.MemberStatus_MEMBER_STATUS_INACTIVE
	case "visitor":
		return membersv1.MemberStatus_MEMBER_STATUS_VISITOR
	default:
		return membersv1.MemberStatus_MEMBER_STATUS_ACTIVE
	}
}
