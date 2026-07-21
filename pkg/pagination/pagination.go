package pagination

import "strings"

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type ListRequest struct {
	Page     int
	PageSize int
	Search   string
}

type Query struct {
	Limit  int
	Offset int
	Search string
}

func (r ListRequest) ToQuery() Query {
	page := r.Page
	if page < 1 {
		page = 1
	}
	pageSize := r.PageSize
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Query{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
		Search: strings.TrimSpace(r.Search),
	}
}

type Page[T any] struct {
	Items    []T
	Total    int
	Page     int
	PageSize int
}

func NewPage[T any](items []T, total int, q Query) Page[T] {
	page := 1
	if q.Limit > 0 {
		page = q.Offset/q.Limit + 1
	}
	return Page[T]{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: q.Limit,
	}
}
