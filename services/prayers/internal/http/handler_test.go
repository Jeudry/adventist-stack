package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/domain"
	prayerhttp "github.com/Jeudry/adventist-stack/services/prayers/internal/http"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/service"
)

type fakeRepo struct {
	created []domain.Prayer
}

func (f *fakeRepo) Create(_ context.Context, p domain.Prayer) (domain.Prayer, error) {
	p.ID = uuid.New()
	f.created = append(f.created, p)
	return p, nil
}

func (f *fakeRepo) Update(_ context.Context, p domain.Prayer) (domain.Prayer, error) { return p, nil }
func (f *fakeRepo) Delete(_ context.Context, _, _ uuid.UUID) error                   { return nil }
func (f *fakeRepo) GetByID(_ context.Context, _ uuid.UUID) (domain.Prayer, error) {
	return domain.Prayer{}, nil
}

func (f *fakeRepo) RetrieveList(_ context.Context, _ pagination.Query) ([]domain.Prayer, error) {
	return nil, nil
}
func (f *fakeRepo) Count(_ context.Context, _ pagination.Query) (int, error) { return 0, nil }

// The real chi stack on purpose: httpx.Identity is chi middleware, so a humatest router would skip
// it and both cases below would pass while proving nothing.
func postPrayer(t *testing.T, repo *fakeRepo, caller string) *httptest.ResponseRecorder {
	t.Helper()

	router := chi.NewMux()
	prayerhttp.NewHandler(service.NewPrayerService(repo)).Register(httpx.NewAPI(router, "prayers-test"))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/prayers",
		strings.NewReader(`{"title":"Petición de prueba","description":"Descripción de la petición.","isAnonymous":true}`))
	req.Header.Set("Content-Type", "application/json")
	if caller != "" {
		req.Header.Set(httpx.UserIDHeader, caller)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

// An anonymous prayer hides the author from other members; it does not make the request itself
// anonymous. The caller is still recorded in created_by.
func TestCreateRecordsTheCallerEvenWhenThePrayerIsAnonymous(t *testing.T) {
	repo := &fakeRepo{}
	actor := uuid.New()

	resp := postPrayer(t, repo, actor.String())

	require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())
	require.Len(t, repo.created, 1)
	assert.Equal(t, actor, repo.created[0].CreatedBy, "created_by must come from X-User-Id")
	assert.Nil(t, repo.created[0].AuthorName, "an anonymous prayer stores no author")
}

func TestCreateWithoutTheHeaderIsRejected(t *testing.T) {
	repo := &fakeRepo{}

	resp := postPrayer(t, repo, "")

	assert.Equal(t, http.StatusUnauthorized, resp.Code, resp.Body.String())
	assert.Empty(t, repo.created, "nothing may be written for a caller that does not exist")
}
