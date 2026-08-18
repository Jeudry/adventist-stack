package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/domain"
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/service"
)

// The service decides; it does not speak SQL. A fake that records the calls is what lets a test
// assert the thing the database cannot show: that an invalid school never reached the repository.
type fakeRepo struct {
	created, updated []domain.SabbathSchool
	deleted          []uuid.UUID
	list             []domain.SabbathSchool
	total            int
}

func (f *fakeRepo) Create(_ context.Context, ss domain.SabbathSchool) (domain.SabbathSchool, error) {
	f.created = append(f.created, ss)
	return ss, nil
}

func (f *fakeRepo) Update(_ context.Context, ss domain.SabbathSchool) (domain.SabbathSchool, error) {
	f.updated = append(f.updated, ss)
	return ss, nil
}

func (f *fakeRepo) Delete(_ context.Context, id, _ uuid.UUID) error {
	f.deleted = append(f.deleted, id)
	return nil
}

func (f *fakeRepo) GetByID(_ context.Context, _ uuid.UUID) (domain.SabbathSchool, error) {
	return domain.SabbathSchool{}, nil
}

func (f *fakeRepo) RetrieveList(_ context.Context, _ pagination.Query) ([]domain.SabbathSchool, error) {
	return f.list, nil
}

func (f *fakeRepo) Count(_ context.Context, _ pagination.Query) (int, error) {
	return f.total, nil
}

func TestCreateRejectsAnInvalidSchoolWithoutTouchingTheRepository(t *testing.T) {
	repo := &fakeRepo{}
	svc := service.NewSabbathSchoolService(repo)

	_, err := svc.Create(t.Context(), domain.SabbathSchool{Status: domain.STATUS_ACTIVE})

	assert.Error(t, err, "creating a school with no name")
	assert.Empty(t, repo.created, "the repository must never see an invalid school")
}

func TestCreateNormalizesBeforeReachingTheRepository(t *testing.T) {
	repo := &fakeRepo{}
	svc := service.NewSabbathSchoolService(repo)

	_, err := svc.Create(t.Context(), domain.SabbathSchool{
		Name:   "   Clase de jóvenes   ",
		Status: domain.STATUS_ACTIVE,
	})

	require.NoError(t, err)
	require.Len(t, repo.created, 1)
	assert.Equal(t, "Clase de jóvenes", repo.created[0].Name, "the name must arrive trimmed")
}

func TestUpdateRejectsTheZeroID(t *testing.T) {
	repo := &fakeRepo{}
	svc := service.NewSabbathSchoolService(repo)

	_, err := svc.Update(t.Context(), domain.SabbathSchool{Name: "Clase", Status: domain.STATUS_ACTIVE})

	assert.ErrorIs(t, err, domain.ErrSabbathSchoolNotFound, "updating with a zero id")
	assert.Empty(t, repo.updated, "the repository must not be asked to update nothing")
}

// The page carries what the query asked for; the total comes from a second query. Deriving it from
// len(items) instead shows up as "2 results" next to a list that has a second page.
func TestRetrieveListTakesTheTotalFromCount(t *testing.T) {
	repo := &fakeRepo{
		list:  []domain.SabbathSchool{{Name: "Una"}, {Name: "Otra"}},
		total: 57,
	}
	svc := service.NewSabbathSchoolService(repo)

	page, err := svc.RetrieveList(t.Context(), pagination.Query{Limit: 2, Offset: 0})

	require.NoError(t, err)
	assert.Len(t, page.Items, 2, "the page carries only what the query asked for")
	assert.Equal(t, 57, page.Total, "the total must come from Count and not from the page length")
}
