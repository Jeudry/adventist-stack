package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/domain"
	sshttp "github.com/Jeudry/adventist-stack/services/sabbath_school/internal/http"
	"github.com/Jeudry/adventist-stack/services/sabbath_school/internal/service"
)

type fakeRepo struct {
	created []domain.SabbathSchool
}

func (f *fakeRepo) Create(_ context.Context, ss domain.SabbathSchool) (domain.SabbathSchool, error) {
	ss.ID = uuid.New()
	f.created = append(f.created, ss)
	return ss, nil
}

func (f *fakeRepo) Update(_ context.Context, ss domain.SabbathSchool) (domain.SabbathSchool, error) {
	return ss, nil
}
func (f *fakeRepo) Delete(_ context.Context, _, _ uuid.UUID) error { return nil }
func (f *fakeRepo) GetByID(_ context.Context, _ uuid.UUID) (domain.SabbathSchool, error) {
	return domain.SabbathSchool{}, nil
}
func (f *fakeRepo) RetrieveList(_ context.Context, _ pagination.Query) ([]domain.SabbathSchool, error) {
	return nil, nil
}
func (f *fakeRepo) Count(_ context.Context, _ pagination.Query) (int, error) { return 0, nil }

func newHandler(repo *fakeRepo) *sshttp.Handler {
	return sshttp.NewHandler(service.NewSabbathSchoolService(repo))
}

// humatest brings its own router, so the chi middleware chain — Identity included — never runs.
// That is what these cases want: binding and validation on their own.
func TestValidation(t *testing.T) {
	cases := map[string]struct {
		method, path string
		body         any
	}{
		"a body without name":      {method: http.MethodPost, path: "/api/v1/sabbath-schools", body: map[string]any{"status": "active"}},
		"a name below the minimum": {method: http.MethodPost, path: "/api/v1/sabbath-schools", body: map[string]any{"name": "ab", "status": "active"}},
		"an id that is not a uuid": {method: http.MethodGet, path: "/api/v1/sabbath-schools/no-soy-un-uuid"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, api := humatest.New(t)
			newHandler(&fakeRepo{}).Register(api)

			var resp *httptest.ResponseRecorder
			if tc.body != nil {
				resp = api.Do(tc.method, tc.path, "Authorization: irrelevant", tc.body)
			} else {
				resp = api.Do(tc.method, tc.path)
			}

			assert.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
		})
	}
}

// This one mounts the real chi stack on purpose: httpx.Identity is chi middleware, so with
// humatest it would never run and both cases below would pass while proving nothing.
func newRealStack(t *testing.T, repo *fakeRepo) *chi.Mux {
	t.Helper()

	router := chi.NewMux()
	newHandler(repo).Register(httpx.NewAPI(router, "sabbath-school-test"))
	return router
}

func postSchool(t *testing.T, router *chi.Mux, caller string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sabbath-schools",
		strings.NewReader(`{"name":"Clase de prueba","status":"active"}`))
	req.Header.Set("Content-Type", "application/json")
	if caller != "" {
		req.Header.Set(httpx.UserIDHeader, caller)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func TestCreateRecordsTheCallerFromTheHeader(t *testing.T) {
	repo := &fakeRepo{}
	actor := uuid.New()

	resp := postSchool(t, newRealStack(t, repo), actor.String())

	require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())
	require.Len(t, repo.created, 1)
	assert.Equal(t, actor, repo.created[0].CreatedBy, "created_by must come from X-User-Id")
}

// Every write sits behind the gateway's Auth middleware, so a request with no caller means it
// skipped the gateway. Storing uuid.Nil would put a user that does not exist in created_by.
func TestCreateWithoutTheHeaderIsRejected(t *testing.T) {
	repo := &fakeRepo{}

	resp := postSchool(t, newRealStack(t, repo), "")

	assert.Equal(t, http.StatusUnauthorized, resp.Code, resp.Body.String())
	assert.Empty(t, repo.created, "nothing may be written for a caller that does not exist")
}
