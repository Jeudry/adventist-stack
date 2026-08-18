package httpx

import (
	"time"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/Jeudry/adventist-stack/pkg/pagination"
)

// DateLayout is the wire format for plain dates. Instants keep RFC3339.
const DateLayout = "2006-01-02"

type BaseVM struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	CreatedBy string  `json:"createdBy"`
	UpdatedBy *string `json:"updatedBy,omitempty"`
}

func ToBaseVM(b entity.Base) BaseVM {
	return BaseVM{
		ID:        b.ID.String(),
		CreatedAt: b.CreatedAt.Format(time.RFC3339),
		UpdatedAt: b.UpdatedAt.Format(time.RFC3339),
	}
}

type PageResponse[T any] struct {
	Items    []T `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func ToPageResponse[D any, V any](page pagination.Page[D], toView func(D) V) PageResponse[V] {
	items := make([]V, len(page.Items))
	for i, item := range page.Items {
		items[i] = toView(item)
	}
	return PageResponse[V]{
		Items:    items,
		Total:    page.Total,
		Page:     page.Page,
		PageSize: page.PageSize,
	}
}

func ParseDate(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	parsed, err := time.Parse(DateLayout, *s)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func FormatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format(DateLayout)
	return &formatted
}
