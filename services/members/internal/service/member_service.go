package service

import (
	"context"
	"fmt"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
	"github.com/google/uuid"
)

type memberRepository interface {
	Create(ctx context.Context, m domain.Member) (domain.Member, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Member, error)
	RetrieveList(ctx context.Context, pq pagination.Query) ([]domain.Member, error)
	Count(ctx context.Context, pq pagination.Query) (int, error)
	Update(ctx context.Context, m domain.Member) (domain.Member, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type MemberService struct {
	repo memberRepository
}

func NewMemberService(repo memberRepository) *MemberService {
	return &MemberService{repo: repo}
}

func (s MemberService) Create(ctx context.Context, member domain.Member) (domain.Member, error) {
	m := member.Normalize()
	if err := m.Validate(); err != nil {
		return domain.Member{}, err
	}

	return s.repo.Create(ctx, m)
}

func (s MemberService) GetByID(ctx context.Context, member domain.Member) (domain.Member, error) {
	return s.repo.GetByID(ctx, member.ID)
}

func (s MemberService) RetrieveList(ctx context.Context, req pagination.ListRequest) (pagination.Page[domain.Member], error) {
	q := req.ToQuery()

	items, err := s.repo.RetrieveList(ctx, q)
	if err != nil {
		return pagination.Page[domain.Member]{}, err
	}
	total, err := s.repo.Count(ctx, q)
	if err != nil {
		return pagination.Page[domain.Member]{}, err
	}

	return pagination.NewPage(items, total, q), nil
}

func (s MemberService) Update(ctx context.Context, member domain.Member) (domain.Member, error) {
	if member.ID == uuid.Nil {
		return domain.Member{}, fmt.Errorf("%w: domain update id is required", domain.ErrorInvalidMember)
	}

	m := member.Normalize()
	if err := m.Validate(); err != nil {
		return domain.Member{}, err
	}

	return s.repo.Update(ctx, m)
}

func (s MemberService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
