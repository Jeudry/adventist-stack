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
	"github.com/Jeudry/adventist-stack/services/members/internal/domain"
	memberhttp "github.com/Jeudry/adventist-stack/services/members/internal/http"
	"github.com/Jeudry/adventist-stack/services/members/internal/service"
)

type fakeRepo struct {
	created []domain.Member
}

func (f *fakeRepo) Create(_ context.Context, m domain.Member) (domain.Member, error) {
	m.ID = uuid.New()
	f.created = append(f.created, m)
	return m, nil
}

func (f *fakeRepo) Update(_ context.Context, m domain.Member) (domain.Member, error) { return m, nil }
func (f *fakeRepo) Delete(_ context.Context, _, _ uuid.UUID) error                   { return nil }
func (f *fakeRepo) GetByID(_ context.Context, _ uuid.UUID) (domain.Member, error) {
	return domain.Member{}, nil
}

func (f *fakeRepo) RetrieveList(_ context.Context, _ pagination.Query) ([]domain.Member, error) {
	return nil, nil
}
func (f *fakeRepo) Count(_ context.Context, _ pagination.Query) (int, error) { return 0, nil }

// The real chi stack on purpose: httpx.Identity is chi middleware, so a humatest router would skip
// it and both cases below would pass while proving nothing.
func postMember(t *testing.T, repo *fakeRepo, caller string) *httptest.ResponseRecorder {
	t.Helper()

	router := chi.NewMux()
	memberhttp.NewHandler(service.NewMemberService(repo)).Register(httpx.NewAPI(router, "members-test"))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/members",
		strings.NewReader(`{"firstName":"Ana","lastName":"Pérez","gender":"F","status":"active"}`))
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

	resp := postMember(t, repo, actor.String())

	require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())
	require.Len(t, repo.created, 1)
	assert.Equal(t, actor, repo.created[0].CreatedBy, "created_by must come from X-User-Id")
}

func TestCreateWithoutTheHeaderIsRejected(t *testing.T) {
	repo := &fakeRepo{}

	resp := postMember(t, repo, "")

	assert.Equal(t, http.StatusUnauthorized, resp.Code, resp.Body.String())
	assert.Empty(t, repo.created, "nothing may be written for a caller that does not exist")
}
