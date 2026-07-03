// Package handlers contiene los handlers HTTP del gateway y los viewmodels
// (DTOs) que se serializan hacia/desde el cliente.
package handlers

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// writeJSON serializa v como JSON con el status dado.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// ErrorResponse es el viewmodel estándar de error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// writeError traduce un error (incluyendo errores gRPC) a una respuesta HTTP.
func writeError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	msg := "error interno"

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

// decodeJSON deserializa el body de la request en v.
func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
