package base

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type BaseVM struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	CreatedBy *string `json:"created_by,omitempty"`
	UpdatedBy *string `json:"updated_by,omitempty"`
}

func ToBaseVM(id string, createdAt, updatedAt *timestamppb.Timestamp) BaseVM {
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
