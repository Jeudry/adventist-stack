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
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/domain"
)

func setup(t *testing.T) (context.Context, *pgxpool.Pool, *SabbathSchoolRepository) {
	t.Helper()

	pool := testdb.Connect(t)
	return t.Context(), pool, NewSabbathSchoolRepository(pool)
}

func newSabbathSchool(t *testing.T, pool *pgxpool.Pool, repo *SabbathSchoolRepository, name string, location *string) domain.SabbathSchool {
	t.Helper()

	now := time.Now().UTC()
	ss := domain.SabbathSchool{
		Name:         name,
		Location:     location,
		TargetMinAge: new(15),
		TargetMaxAge: new(25),
		Status:       domain.STATUS_ACTIVE,
	}
	ss.ID = uuid.New()
	ss.CreatedAt, ss.UpdatedAt = now, now
	ss.CreatedBy = uuid.New()

	created, err := repo.Create(t.Context(), ss)
	require.NoError(t, err, "create %q", name)

	t.Cleanup(func() { testdb.Exec(t, pool, "DELETE FROM sabbath_school WHERE id = $1", created.ID) })
	return created
}

func rowIncludingDeleted(t *testing.T, pool *pgxpool.Pool, id uuid.UUID) domain.SabbathSchool {
	t.Helper()

	found, err := scanSabbathSchool(pool.QueryRow(t.Context(), selectSabbathSchoolIncludingDeleted, pgx.StrictNamedArgs{"id": id}))
	require.NoError(t, err, "read row including deleted")
	return found
}

func TestSabbathSchoolRepositoryRoundTrip(t *testing.T) {
	ctx, pool, repo := setup(t)

	created := newSabbathSchool(t, pool, repo, "Clase de jóvenes", new("Salón 3"))

	assert.NotEqual(t, uuid.Nil, created.ID, "id")
	assert.Equal(t, "Salón 3", ptr.Deref(created.Location), "location")
	assert.Equal(t, 15, ptr.Deref(created.TargetMinAge), "target_min_age")
	assert.Equal(t, 25, ptr.Deref(created.TargetMaxAge), "target_max_age")
	assert.Equal(t, domain.STATUS_ACTIVE, created.Status, "status")
	assert.Equal(t, uuid.NullUUID{}, created.TeacherID, "teacher_id")

	got, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err, "get by id")
	assert.Equal(t, "Clase de jóvenes", got.Name, "name")
}

func TestSabbathSchoolRepositoryEmptyLocationStoresNull(t *testing.T) {
	ctx, pool, repo := setup(t)

	created := newSabbathSchool(t, pool, repo, "Clase sin salón", nil)
	assert.Nil(t, created.Location, "returned: location: got %q", ptr.Deref(created.Location))

	reread, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err, "get by id")
	assert.Nil(t, reread.Location, "reread: location: got %q", ptr.Deref(reread.Location))
}

func TestSabbathSchoolRepositoryGetByIDNotFound(t *testing.T) {
	ctx, _, repo := setup(t)

	_, err := repo.GetByID(ctx, uuid.New())
	assert.ErrorIs(t, err, domain.ErrSabbathSchoolNotFound, "reading a row that never existed")
}

func TestSabbathSchoolRepositoryUpdateWritesEveryField(t *testing.T) {
	ctx, pool, repo := setup(t)

	created := newSabbathSchool(t, pool, repo, "Antes de editar", new("Salón 3"))
	editor := uuid.New()

	edited := created
	edited.Name = "Despues de editar"
	edited.Location = new("Salon 9")
	edited.TargetMinAge = new(7)
	edited.TargetMaxAge = new(12)
	edited.Status = domain.STATUS_INACTIVE
	edited.UpdatedBy = &editor

	updated, err := repo.Update(ctx, edited)
	require.NoError(t, err, "update")

	reread, err := repo.GetByID(ctx, updated.ID)
	require.NoError(t, err, "get by id")

	// Every field takes a value the row did not already hold, and the two ints differ from each
	// other: a swapped pair of arguments is invisible when both sides carry the same number.
	assertMatches := func(t *testing.T, got domain.SabbathSchool) {
		t.Helper()

		assert.Equal(t, "Despues de editar", got.Name, "name")
		assert.Equal(t, "Salon 9", ptr.Deref(got.Location), "location")
		assert.Equal(t, 7, ptr.Deref(got.TargetMinAge), "target_min_age")
		assert.Equal(t, 12, ptr.Deref(got.TargetMaxAge), "target_max_age")
		assert.Equal(t, domain.STATUS_INACTIVE, got.Status, "status")
		assert.Equal(t, editor, ptr.Deref(got.UpdatedBy), "updated_by")
		assert.Greater(t, got.UpdatedAt, created.UpdatedAt, "updated_at")
		assert.WithinDuration(t, created.CreatedAt, got.CreatedAt, 0, "created_at must stay untouched")
	}

	t.Run("returned", func(t *testing.T) { assertMatches(t, updated) })
	t.Run("reread", func(t *testing.T) {
		assertMatches(t, reread)
		assert.WithinDuration(t, updated.UpdatedAt, reread.UpdatedAt, 0, "updated_at must survive the read")
	})
}

func TestSabbathSchoolRepositoryUpdateClearsOptionalFields(t *testing.T) {
	ctx, pool, repo := setup(t)

	created := newSabbathSchool(t, pool, repo, "Clase con todo", new("Salón 3"))

	cleared := created
	cleared.Location = nil
	cleared.TargetMinAge = nil
	cleared.TargetMaxAge = nil

	updated, err := repo.Update(ctx, cleared)
	require.NoError(t, err, "update")

	reread, err := repo.GetByID(ctx, updated.ID)
	require.NoError(t, err, "get by id")

	assertCleared := func(t *testing.T, got domain.SabbathSchool) {
		t.Helper()

		assert.Nil(t, got.Location, "location: got %q", ptr.Deref(got.Location))
		assert.Nil(t, got.TargetMinAge, "target_min_age: got %d", ptr.Deref(got.TargetMinAge))
		assert.Nil(t, got.TargetMaxAge, "target_max_age: got %d", ptr.Deref(got.TargetMaxAge))
	}

	t.Run("returned", func(t *testing.T) { assertCleared(t, updated) })
	t.Run("reread", func(t *testing.T) { assertCleared(t, reread) })
}

func TestSabbathSchoolRepositoryDeleteIsSoft(t *testing.T) {
	ctx, pool, repo := setup(t)

	created := newSabbathSchool(t, pool, repo, "Clase a borrar", new("Salón 3"))
	deleter := uuid.New()
	require.NoError(t, repo.Delete(ctx, created.ID, deleter), "delete")

	_, err := repo.GetByID(ctx, created.ID)
	assert.ErrorIs(t, err, domain.ErrSabbathSchoolNotFound, "reading a deleted row")
	assert.ErrorIs(t, repo.Delete(ctx, created.ID, deleter), domain.ErrSabbathSchoolNotFound, "deleting twice")

	// The row surviving with the audit columns filled is what separates a soft delete from a
	// real one, and no query the repository exposes can see it: they all filter deleted_at.
	row := rowIncludingDeleted(t, pool, created.ID)
	assert.NotNil(t, row.DeletedAt, "deleted_at")
	assert.Equal(t, deleter, ptr.Deref(row.DeletedBy), "deleted_by")
}

func TestSabbathSchoolRepositoryUpdateOnDeletedFails(t *testing.T) {
	ctx, pool, repo := setup(t)

	created := newSabbathSchool(t, pool, repo, "Clase borrada", new("Salón 3"))
	require.NoError(t, repo.Delete(ctx, created.ID, uuid.New()), "delete")

	edited := created
	edited.Name = "No debería escribirse"

	_, err := repo.Update(ctx, edited)
	assert.ErrorIs(t, err, domain.ErrSabbathSchoolNotFound, "updating a deleted row")

	// The right error does not prove nothing was written: an UPDATE missing its deleted_at guard
	// could write the row and still report no rows back.
	assert.Equal(t, "Clase borrada", rowIncludingDeleted(t, pool, created.ID).Name, "name must stay untouched")
}

func TestSabbathSchoolRepositoryPagination(t *testing.T) {
	_, pool, repo := setup(t)

	paginationtest.Run(t, paginationtest.Subject[domain.SabbathSchool]{
		List:  repo.RetrieveList,
		Count: repo.Count,
		Create: func(t *testing.T, name string) domain.SabbathSchool {
			return newSabbathSchool(t, pool, repo, name, new("Salón 3"))
		},
		Delete: func(t *testing.T, ss domain.SabbathSchool) {
			require.NoError(t, repo.Delete(t.Context(), ss.ID, uuid.New()), "delete")
		},
		ID: func(ss domain.SabbathSchool) uuid.UUID { return ss.ID },
	})
}

func TestSabbathSchoolRepositoryStoresTheTeacher(t *testing.T) {
	ctx, pool, repo := setup(t)

	teacherId := uuid.New()
	ss := domain.SabbathSchool{
		Name:      "Clase con maestro",
		TeacherID: uuid.NullUUID{UUID: teacherId, Valid: true},
		Status:    domain.STATUS_ACTIVE,
	}
	ss.ID, ss.CreatedBy = uuid.New(), uuid.New()

	created, err := repo.Create(ctx, ss)
	require.NoError(t, err, "create")
	t.Cleanup(func() { testdb.Exec(t, pool, "DELETE FROM sabbath_school WHERE id = $1", created.ID) })

	reread, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err, "get by id")

	// The uuid.NullUUID crossing has two branches in each direction and every other test only
	// exercises the nil one, so an inverted check in either mapper would go unnoticed.
	assertTeacher := func(t *testing.T, got domain.SabbathSchool) {
		t.Helper()

		assert.True(t, got.TeacherID.Valid, "teacher_id must not be NULL")
		assert.Equal(t, teacherId, got.TeacherID.UUID, "teacher_id")
	}

	t.Run("returned", func(t *testing.T) { assertTeacher(t, created) })
	t.Run("reread", func(t *testing.T) { assertTeacher(t, reread) })

	t.Run("clearing it reaches the column", func(t *testing.T) {
		cleared := reread
		cleared.TeacherID = uuid.NullUUID{}

		updated, err := repo.Update(ctx, cleared)
		require.NoError(t, err, "update")
		assert.Equal(t, uuid.NullUUID{}, updated.TeacherID, "teacher_id after clearing")
	})
}

func TestSabbathSchoolRepositoryListSearchesEveryColumn(t *testing.T) {
	ctx, pool, repo := setup(t)

	// The list query offers three OR branches and every other test searches by name, so the name
	// branch is answering for all of them: a wrong column or a broken cast in the other two is
	// invisible today. That is the shape of the ILIKE-on-an-integer bug this schema already had.
	location := "Salon " + uuid.NewString()
	created := newSabbathSchool(t, pool, repo, "Clase buscada por otra columna", &location)

	byLocation, err := repo.RetrieveList(ctx, pagination.Query{Search: location, Limit: 10})
	require.NoError(t, err, "search by location")
	require.Len(t, byLocation, 1, "search by location")
	assert.Equal(t, created.ID, byLocation[0].ID, "search by location")

	// Ordering is created_at DESC and this row is the newest, so a full page reaches it.
	byStatus, err := repo.RetrieveList(ctx, pagination.Query{Search: domain.STATUS_ACTIVE.String(), Limit: 100})
	require.NoError(t, err, "search by status")
	assert.Contains(t, idsOf(byStatus), created.ID, "search by status must reach the status column")
}

func idsOf(rows []domain.SabbathSchool) []uuid.UUID {
	ids := make([]uuid.UUID, len(rows))
	for i, row := range rows {
		ids[i] = row.ID
	}
	return ids
}
