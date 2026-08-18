//go:build integration

package repository

import (
	"context"
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
	"github.com/Jeudry/adventist-stack/services/prayers/internal/domain"
)

func setup(t *testing.T) (context.Context, *pgxpool.Pool, *PrayerRepository) {
	t.Helper()

	pool := testdb.Connect(t)
	return t.Context(), pool, NewPrayerRepository(pool)
}

// A nil author is an anonymous prayer; there is no second field to keep in step.
func newPrayer(title string, author *string) domain.Prayer {
	now := time.Now().UTC()
	p := domain.Prayer{
		Title:       title,
		Description: "Descripción de la petición de prueba.",
		AuthorName:  author,
		Status:      domain.STATUS_PENDING,
	}
	p.ID = uuid.New()
	p.CreatedAt = now
	p.UpdatedAt = now
	p.CreatedBy = uuid.New()
	return p
}

func rowIncludingDeleted(t *testing.T, pool *pgxpool.Pool, id uuid.UUID) domain.Prayer {
	t.Helper()

	found, err := scanPrayer(pool.QueryRow(t.Context(), selectPrayerIncludingDeleted, pgx.StrictNamedArgs{"id": id}))
	require.NoError(t, err, "read row including deleted")
	return found
}

// An anonymous prayer stores author_name as NULL. Reading it back used to panic:
// the old mapper dereferenced the pointer without a nil check.
func TestPrayerRepositoryAnonymousRoundTrip(t *testing.T) {
	ctx, pool, repo := setup(t)

	created, err := repo.Create(ctx, newPrayer("Petición anónima", nil))
	require.NoError(t, err, "create")
	t.Cleanup(func() { testdb.Exec(t, pool, "DELETE FROM prayers WHERE id = $1", created.ID) })

	assert.Nil(t, created.AuthorName, "returned: author_name: got %q", ptr.Deref(created.AuthorName))

	got, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err, "get by id")
	assert.Nil(t, got.AuthorName, "reread: author_name: got %q", ptr.Deref(got.AuthorName))
}

func TestPrayerRepositoryCRUD(t *testing.T) {
	ctx, pool, repo := setup(t)

	created, err := repo.Create(ctx, newPrayer("Petición con autor", new("Hermano Luis")))
	require.NoError(t, err, "create")
	t.Cleanup(func() { testdb.Exec(t, pool, "DELETE FROM prayers WHERE id = $1", created.ID) })

	assert.Equal(t, domain.STATUS_PENDING, created.Status, "status")
	assert.Equal(t, "Hermano Luis", ptr.Deref(created.AuthorName), "author_name")

	query := pagination.Query{Search: "Petición con autor", Limit: 10}

	list, err := repo.RetrieveList(ctx, query)
	require.NoError(t, err, "retrieve list")
	assert.NotEmpty(t, list, "retrieve list: the search found nothing")

	count, err := repo.Count(ctx, pagination.Query{Search: query.Search})
	require.NoError(t, err, "count")
	assert.Equal(t, len(list), count, "count must agree with the list length")

	created.Title = "Petición actualizada"
	updated, err := repo.Update(ctx, created)
	require.NoError(t, err, "update")
	assert.Equal(t, "Petición actualizada", updated.Title, "title")

	deleter := uuid.New()
	require.NoError(t, repo.Delete(ctx, created.ID, deleter), "delete")

	_, err = repo.GetByID(ctx, created.ID)
	assert.ErrorIs(t, err, domain.ErrPrayerNotFound, "reading a soft-deleted prayer")
	assert.ErrorIs(t, repo.Delete(ctx, created.ID, deleter), domain.ErrPrayerNotFound, "deleting twice")

	row := rowIncludingDeleted(t, pool, created.ID)
	assert.NotNil(t, row.DeletedAt, "deleted_at")
	assert.Equal(t, deleter, ptr.Deref(row.DeletedBy), "deleted_by")
}

func TestPrayerRepositoryGetByIDNotFound(t *testing.T) {
	ctx, _, repo := setup(t)

	_, err := repo.GetByID(ctx, uuid.New())
	assert.ErrorIs(t, err, domain.ErrPrayerNotFound, "reading a row that never existed")
}

func createPrayer(t *testing.T, pool *pgxpool.Pool, repo *PrayerRepository, title string) domain.Prayer {
	t.Helper()

	created, err := repo.Create(t.Context(), newPrayer(title, new("Hermano Luis")))
	require.NoError(t, err, "create %q", title)
	t.Cleanup(func() { testdb.Exec(t, pool, "DELETE FROM prayers WHERE id = $1", created.ID) })
	return created
}

func TestPrayerRepositoryPagination(t *testing.T) {
	_, pool, repo := setup(t)

	paginationtest.Run(t, paginationtest.Subject[domain.Prayer]{
		List:  repo.RetrieveList,
		Count: repo.Count,
		Create: func(t *testing.T, name string) domain.Prayer {
			return createPrayer(t, pool, repo, name)
		},
		Delete: func(t *testing.T, pr domain.Prayer) {
			require.NoError(t, repo.Delete(t.Context(), pr.ID, uuid.New()), "delete")
		},
		ID: func(pr domain.Prayer) uuid.UUID { return pr.ID },
	})
}
