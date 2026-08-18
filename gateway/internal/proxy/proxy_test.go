package proxy_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jeudry/adventist-stack/gateway/internal/proxy"
	"github.com/Jeudry/adventist-stack/pkg/httpx"
	"github.com/Jeudry/adventist-stack/pkg/middleware"
)

type received struct {
	path   string
	userID string
	role   string
}

func upstreamRecording(t *testing.T, got *received) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.path = r.URL.Path
		got.userID = r.Header.Get(httpx.UserIDHeader)
		got.role = r.Header.Get(httpx.RoleHeader)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	return server
}

// Every service declares the same public path the client called, so rewriting
// it here would send the request somewhere the service never registered.
func TestTo_ForwardsThePathUntouched(t *testing.T) {
	var got received
	upstream := upstreamRecording(t, &got)

	handler, err := proxy.To(upstream.URL, slog.Default())
	require.NoError(t, err, "build proxy")

	request := httptest.NewRequest(http.MethodGet, "/api/v1/members/abc-123", nil)
	handler.ServeHTTP(httptest.NewRecorder(), request)

	assert.Equal(t, "/api/v1/members/abc-123", got.path, "the path the upstream saw")
}

func TestTo_AssertsTheAuthenticatedCaller(t *testing.T) {
	var got received
	upstream := upstreamRecording(t, &got)

	handler, err := proxy.To(upstream.URL, slog.Default())
	require.NoError(t, err, "build proxy")

	request := httptest.NewRequest(http.MethodGet, "/api/v1/members/", nil)
	handler.ServeHTTP(httptest.NewRecorder(), request.WithContext(
		middleware.WithCaller(context.Background(), "user-42", "admin"),
	))

	assert.Equal(t, "user-42", got.userID, httpx.UserIDHeader)
	assert.Equal(t, "admin", got.role, httpx.RoleHeader)
}

// A client that sends the identity headers itself must not be believed: they are
// the gateway's word about who the caller is, and the services trust them blindly.
func TestTo_DropsIdentityHeadersSentByTheClient(t *testing.T) {
	var got received
	upstream := upstreamRecording(t, &got)

	handler, err := proxy.To(upstream.URL, slog.Default())
	require.NoError(t, err, "build proxy")

	request := httptest.NewRequest(http.MethodGet, "/api/v1/members/", nil)
	request.Header.Set(httpx.UserIDHeader, "somebody-else")
	request.Header.Set(httpx.RoleHeader, "admin")
	handler.ServeHTTP(httptest.NewRecorder(), request)

	assert.Empty(t, got.userID, "the upstream must not see a forged %s", httpx.UserIDHeader)
	assert.Empty(t, got.role, "the upstream must not see a forged %s", httpx.RoleHeader)
}

func TestTo_AnswersWhenTheServiceIsDown(t *testing.T) {
	handler, err := proxy.To("http://127.0.0.1:1", slog.New(slog.DiscardHandler))
	require.NoError(t, err, "build proxy")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/members/", nil))

	assert.Equal(t, http.StatusBadGateway, recorder.Code, "status")
	assert.NotEmpty(t, recorder.Body.String(), "a 502 must still carry a JSON body")
}
