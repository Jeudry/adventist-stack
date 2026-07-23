package handlers

import (
	"strings"
	"time"

	membersv1 "github.com/Jeudry/adventist-stack/gen/members/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
