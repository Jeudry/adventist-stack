package handlers

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	msg := "internal error"

	if st, ok := status.FromError(err); ok {
		msg = st.Message()
		switch st.Code() {
		case codes.InvalidArgument:
			code = http.StatusBadRequest
		case codes.Unauthenticated:
			code = http.StatusUnauthorized
		case codes.NotFound:
			code = http.StatusNotFound
		case codes.AlreadyExists:
			code = http.StatusConflict
		case codes.PermissionDenied:
			code = http.StatusForbidden
		}
	}

	writeJSON(w, code, ErrorResponse{Error: msg})
}

func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
