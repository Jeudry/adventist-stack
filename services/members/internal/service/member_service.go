package service

import (
	"context"

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

func (s MemberService) Create(ctx context.Context, domain domain.Member) {
	d := normalize(req)
	if err != validate(p); err != nil {
		return domain.Member{}
	}

	return s.repo.Create()
}

func (s MemberService) GetByID(ctx context.Context, member domain.Member) (domain.Member, error) {
	return s.repo.GetByID(ctx, member.Id)
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
