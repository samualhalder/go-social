package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (p PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	query := r.URL.Query()

	limit := query.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return p, nil
		}
		p.Limit = l
	}
	offset := query.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return p, nil
		}
		p.Offset = o
	}
	sort := query.Get("sort")

	if sort != "" {
		p.Sort = sort
	}
	return p, nil
}
