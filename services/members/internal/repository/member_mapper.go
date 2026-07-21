package repository

import (
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/members/internal/db"
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
)

func toDomain(m db.Member) domain.Member {
	return domain.Member{
		Id:          m.ID,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Email:       m.Email,
		Phone:       m.Phone,
		Gender:      m.Gender,
		Address:     m.Address,
		BirthDate:   m.BirthDate,
		BaptismDate: m.BaptismDate,
		Status:      domain.Status(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toCreateParams(m domain.Member) db.CreateMemberParams {
	return db.CreateMemberParams{
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Email:       m.Email,
		Phone:       m.Phone,
		Gender:      m.Gender,
		Address:     m.Address,
		BirthDate:   m.BirthDate,
		BaptismDate: m.BaptismDate,
		Status:      string(m.Status),
	}
}

func toUpdateParams(m domain.Member) db.UpdateMemberParams {
	return db.UpdateMemberParams{
		ID:          m.Id,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Email:       m.Email,
		Phone:       m.Phone,
		Gender:      m.Gender,
		Address:     m.Address,
		BirthDate:   m.BirthDate,
		BaptismDate: m.BaptismDate,
		Status:      string(m.Status),
	}
}

func toMemberListParams(q pagination.Query) db.ListMembersParams {
	return db.ListMembersParams{
		Search:    q.Search,
		RowLimit:  int32(q.Limit),
		RowOffset: int32(q.Offset),
	}
}
