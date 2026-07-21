package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/auth/internal/db"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{q: db.New(pool)}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	created, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Email:        u.Email.String(),
		Name:         u.Name,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("repository: create user: %w", err)
	}
	return toDomain(created)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email vo.Email) (domain.User, error) {
	u, err := r.q.GetUserByEmail(ctx, email.String())
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("repository: find by email: %w", err)
	}
	return toDomain(u)
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	exists, err := r.q.ExistsUserByEmail(ctx, email.String())
	if err != nil {
		return false, fmt.Errorf("repository: exists by email: %w", err)
	}
	return exists, nil
}
