package repository

import (
	"fmt"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/members/internal/db"
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
)

func toDomain(m db.Member) (domain.Member, error) {
	email, err := vo.NewOptionalEmail(m.Email)
	if err != nil {
		return domain.Member{}, fmt.Errorf("repository: rehydrate email: %w", err)
	}
	phone, err := vo.NewOptionalPhone(m.Phone)
	if err != nil {
		return domain.Member{}, fmt.Errorf("repository: rehydrate phone: %w", err)
	}
	return domain.Member{
		Base: entity.Base{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: m.DeletedAt,
			CreatedBy: m.CreatedBy,
			UpdatedBy: m.UpdatedBy,
			DeletedBy: m.DeletedBy,
		},
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Email:       email,
		Phone:       phone,
		Gender:      genderFromDB(m.Gender),
		Address:     m.Address,
		BirthDate:   m.BirthDate,
		BaptismDate: m.BaptismDate,
		Status:      statusFromDB(m.Status),
	}, nil
}

func statusFromDB(s int16) domain.Status {
	switch s {
	case 2:
		return domain.StatusInactive
	case 3:
		return domain.StatusVisitor
	default:
		return domain.StatusActive
	}
}

func genderFromDB(g int16) domain.Gender {
	if g == 2 {
		return domain.GenderFemale
	}
	return domain.GenderMale
}

func toCreateParams(m domain.Member) db.CreateMemberParams {
	return db.CreateMemberParams{
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Email:       m.Email.Ptr(),
		Phone:       m.Phone.Ptr(),
		Gender:      genderToDB(m.Gender),
		Address:     m.Address,
		BirthDate:   m.BirthDate,
		BaptismDate: m.BaptismDate,
		Status:      statusToDB(m.Status),
	}
}

func toUpdateParams(m domain.Member) db.UpdateMemberParams {
	return db.UpdateMemberParams{
		ID:          m.ID,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Email:       m.Email.Ptr(),
		Phone:       m.Phone.Ptr(),
		Gender:      genderToDB(m.Gender),
		Address:     m.Address,
		BirthDate:   m.BirthDate,
		BaptismDate: m.BaptismDate,
		Status:      statusToDB(m.Status),
	}
}

func statusToDB(s domain.Status) int16 {
	switch s {
	case domain.StatusInactive:
		return 2
	case domain.StatusVisitor:
		return 3
	default:
		return 1
	}
}

func genderToDB(g domain.Gender) int16 {
	if g == domain.GenderFemale {
		return 2
	}
	return 1
}

func toMemberListParams(q pagination.Query) db.ListMembersParams {
	return db.ListMembersParams{
		Search:    q.Search,
		RowLimit:  int32(q.Limit),
		RowOffset: int32(q.Offset),
	}
}
