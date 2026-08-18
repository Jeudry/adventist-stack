//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/pkg/paginationtest"
	"github.com/Jeudry/adventist-stack/pkg/ptr"
	"github.com/Jeudry/adventist-stack/pkg/testdb"
	"github.com/Jeudry/adventist-stack/pkg/vo"
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
)

func setup(t *testing.T) (context.Context, *pgxpool.Pool, *MemberRepository) {
	t.Helper()

	pool := testdb.Connect(t)
	return t.Context(), pool, NewMemberRepository(pool)
}

func newMember(t *testing.T, firstName string) domain.Member {
	t.Helper()

	email, err := vo.NewOptionalEmail(new(fmt.Sprintf("it-%s@iglesia.org", uuid.NewString())))
	require.NoError(t, err, "build email")
	phone, err := vo.NewOptionalPhone(new("+50761234567"))
	require.NoError(t, err, "build phone")
	birth := time.Date(1990, 5, 4, 0, 0, 0, 0, time.UTC)

	m := domain.Member{
		FirstName: firstName,
		LastName:  "Pérez",
		Email:     email,
		Phone:     phone,
		Gender:    domain.GenderFemale,
		Address:   new("Calle 50, Ciudad de Panamá"),
		BirthDate: &birth,
		Status:    domain.StatusActive,
	}
	m.CreatedBy = uuid.New()
	return m
}

func rowIncludingDeleted(t *testing.T, pool *pgxpool.Pool, id uuid.UUID) domain.Member {
	t.Helper()

	found, err := scanMember(pool.QueryRow(t.Context(), selectMemberIncludingDeleted, pgx.StrictNamedArgs{"id": id}))
	require.NoError(t, err, "read row including deleted")
	return found
}

func TestMemberRepositoryRoundTrip(t *testing.T) {
	ctx, pool, repo := setup(t)

	created, err := repo.Create(ctx, newMember(t, "Ana"))
	require.NoError(t, err, "create")
	t.Cleanup(func() { testdb.Exec(t, pool, "DELETE FROM members WHERE id = $1", created.ID) })

	// gender and status are enums, and address sits last in the column list because a later
	// migration added it — both are easy to mis-order.
	assert.Equal(t, domain.GenderFemale, created.Gender, "gender")
	assert.Equal(t, domain.StatusActive, created.Status, "status")
	assert.Equal(t, "Calle 50, Ciudad de Panamá", ptr.Deref(created.Address), "address")
	require.NotNil(t, created.BirthDate, "birth_date")
	assert.Equal(t, 1990, created.BirthDate.Year(), "birth_date year")
	assert.Nil(t, created.BaptismDate, "baptism_date")

	got, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err, "get by id")
	assert.Equal(t, created.Email.String(), got.Email.String(), "email")
	assert.Equal(t, created.Phone.String(), got.Phone.String(), "phone")
}

func TestMemberRepositoryListCountUpdateDelete(t *testing.T) {
	ctx, pool, repo := setup(t)

	created, err := repo.Create(ctx, newMember(t, "Bernarda"))
	require.NoError(t, err, "create")
	t.Cleanup(func() { testdb.Exec(t, pool, "DELETE FROM members WHERE id = $1", created.ID) })

	list, err := repo.RetrieveList(ctx, pagination.Query{Search: "Bernarda", Limit: 10})
	require.NoError(t, err, "retrieve list")
	require.Len(t, list, 1, "retrieve list")
	assert.Equal(t, created.ID, list[0].ID, "the row the search returned")

	count, err := repo.Count(ctx, pagination.Query{Search: "Bernarda"})
	require.NoError(t, err, "count")
	assert.Equal(t, 1, count, "count")

	created.LastName = "Gómez"
	updated, err := repo.Update(ctx, created)
	require.NoError(t, err, "update")
	assert.Equal(t, "Gómez", updated.LastName, "last_name")

	deleter := uuid.New()
	require.NoError(t, repo.Delete(ctx, created.ID, deleter), "delete")

	_, err = repo.GetByID(ctx, created.ID)
	assert.ErrorIs(t, err, domain.ErrMemberNotFound, "reading a soft-deleted member")
	assert.ErrorIs(t, repo.Delete(ctx, created.ID, deleter), domain.ErrMemberNotFound, "deleting twice")

	row := rowIncludingDeleted(t, pool, created.ID)
	assert.NotNil(t, row.DeletedAt, "deleted_at")
	assert.Equal(t, deleter, ptr.Deref(row.DeletedBy), "deleted_by")
}

func createMember(t *testing.T, pool *pgxpool.Pool, repo *MemberRepository, firstName string) domain.Member {
	t.Helper()

	created, err := repo.Create(t.Context(), newMember(t, firstName))
	require.NoError(t, err, "create %q", firstName)
	t.Cleanup(func() { testdb.Exec(t, pool, "DELETE FROM members WHERE id = $1", created.ID) })
	return created
}

func TestMemberRepositoryPagination(t *testing.T) {
	_, pool, repo := setup(t)

	paginationtest.Run(t, paginationtest.Subject[domain.Member]{
		List:  repo.RetrieveList,
		Count: repo.Count,
		Create: func(t *testing.T, name string) domain.Member {
			return createMember(t, pool, repo, name)
		},
		Delete: func(t *testing.T, m domain.Member) {
			require.NoError(t, repo.Delete(t.Context(), m.ID, uuid.New()), "delete")
		},
		ID: func(m domain.Member) uuid.UUID { return m.ID },
	})
}
