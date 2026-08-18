package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
)

// Every statement names its parameters and pgx.StrictNamedArgs refuses to run when one is missing
// or spelled wrong. Positional placeholders fail silently instead: two arguments of the same type
// swapped by mistake still run, and the wrong column is written.
const (
	userColumns = `id, email, name, password_hash, role,
	               created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`

	insertUser = `
INSERT INTO users (email, name, password_hash, role, created_by)
VALUES (@email, @name, @password_hash, @role, @created_by)
RETURNING ` + userColumns

	selectUserByEmail = `
SELECT ` + userColumns + `
FROM users
WHERE email = @email AND deleted_at IS NULL`

	existsUserByEmail = `
SELECT EXISTS (SELECT 1 FROM users WHERE email = @email AND deleted_at IS NULL)`
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	row := r.pool.QueryRow(ctx, insertUser, pgx.StrictNamedArgs{
		"email":         u.Email.String(),
		"name":          u.Name,
		"password_hash": u.PasswordHash,
		"role":          u.Role.String(),
		"created_by":    u.CreatedBy,
	})

	created, err := scanUser(row)
	if err != nil {
		return domain.User{}, fmt.Errorf("repository: create user: %w", err)
	}
	return created, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email vo.Email) (domain.User, error) {
	row := r.pool.QueryRow(ctx, selectUserByEmail, pgx.StrictNamedArgs{"email": email.String()})

	found, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("repository: find by email: %w", err)
	}
	return found, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	var exists bool
	if err := r.pool.QueryRow(ctx, existsUserByEmail, pgx.StrictNamedArgs{"email": email.String()}).Scan(&exists); err != nil {
		return false, fmt.Errorf("repository: exists by email: %w", err)
	}
	return exists, nil
}
