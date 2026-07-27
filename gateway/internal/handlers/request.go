package handlers

import (
	"net/http"
	"strconv"
)

const (
	pageParam     = "page"
	pageSizeParam = "page_size"
	searchParam   = "search"
)

type listQuery struct {
	Page     int32
	PageSize int32
	Search   string
}

// parseListQuery reads the pagination inputs of a list endpoint. Values are forwarded
// as they come: defaults, bounds and trimming belong to pkg/pagination on the service
// side, so applying them here too would fork that policy.
func parseListQuery(r *http.Request) listQuery {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get(pageParam))
	pageSize, _ := strconv.Atoi(query.Get(pageSizeParam))

	return listQuery{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   query.Get(searchParam),
	}
}
