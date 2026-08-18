package httpx

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey int

const userIDKey contextKey = iota

// Identity lifts the caller the gateway authenticated out of UserIDHeader and into the request
// context, which is the only thing a huma handler receives.
//
// A missing or malformed header is not rejected: public routes exist (registering, an anonymous
// prayer) and they reach the same handlers. Those requests simply have no caller, and UserID
// reports uuid.Nil for them.
func Identity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.Header.Get(UserIDHeader))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey, id)))
	})
}

// UserID returns the caller behind the request, or uuid.Nil when it carried no identity.
func UserID(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(userIDKey).(uuid.UUID)
	return id
}
