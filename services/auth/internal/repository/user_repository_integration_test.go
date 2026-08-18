//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jeudry/adventist-stack/pkg/testdb"
	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
)

func setup(t *testing.T) (context.Context, *pgxpool.Pool, *UserRepository) {
	t.Helper()

	pool := testdb.Connect(t)
	return t.Context(), pool, NewUserRepository(pool)
}

// The email carries a timestamp because the column is unique and the table outlives the test.
func uniqueEmail(t *testing.T, prefix string) vo.Email {
	t.Helper()

	email, err := vo.NewEmail(fmt.Sprintf("%s-%d@iglesia.org", prefix, time.Now().UnixNano()))
	require.NoError(t, err, "build email")
	return email
}

func TestUserRepositoryRoundTrip(t *testing.T) {
	ctx, pool, repo := setup(t)

	actor := uuid.New()
	u, err := domain.NewUser(uniqueEmail(t, "it").String(), "Integration Test", domain.RoleAdmin)
	require.NoError(t, err, "build user")
	require.NoError(t, u.SetPassword("supersecret"), "set password")
	u.CreatedBy = actor

	created, err := repo.Create(ctx, u)
	require.NoError(t, err, "create")
	t.Cleanup(func() { testdb.Exec(t, pool, "DELETE FROM users WHERE id = $1", created.ID) })

	assert.NotEqual(t, uuid.Nil, created.ID, "id")
	assert.Equal(t, domain.RoleAdmin, created.Role, "role")
	assert.Equal(t, actor, created.CreatedBy, "created_by")
	assert.False(t, created.CreatedAt.IsZero(), "created_at was not populated by the database")
	assert.Nil(t, created.DeletedAt, "deleted_at")

	found, err := repo.FindByEmail(ctx, created.Email)
	require.NoError(t, err, "find by email")
	assert.Equal(t, created.ID, found.ID, "id")
	assert.Equal(t, created.Email.String(), found.Email.String(), "email")
	assert.NoError(t, found.Authenticate("supersecret"), "the password hash must survive the round-trip")

	exists, err := repo.ExistsByEmail(ctx, created.Email)
	require.NoError(t, err, "exists by email")
	assert.True(t, exists, "exists by email")
}

func TestUserRepositoryFindByEmailNotFound(t *testing.T) {
	ctx, _, repo := setup(t)

	_, err := repo.FindByEmail(ctx, uniqueEmail(t, "nobody"))
	assert.ErrorIs(t, err, domain.ErrUserNotFound, "looking up an email nobody registered")
}
