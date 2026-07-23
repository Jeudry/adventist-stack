package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type BaseVM struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	CreatedBy *string `json:"created_by,omitempty"`
	UpdatedBy *string `json:"updated_by,omitempty"`
}

func toBaseVM(id string, createdAt, updatedAt *timestamppb.Timestamp) BaseVM {
	var createdStr, updatedStr string
	if createdAt != nil {
		createdStr = createdAt.AsTime().Format(time.RFC3339)
	}
	if updatedAt != nil {
		updatedStr = updatedAt.AsTime().Format(time.RFC3339)
	}
	return BaseVM{
		ID:        id,
		CreatedAt: createdStr,
		UpdatedAt: updatedStr,
	}
}

type PageResponse[T any] struct {
	Items    []T   `json:"items"`
	Total    int32 `json:"total"`
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}
