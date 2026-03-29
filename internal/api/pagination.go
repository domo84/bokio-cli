package api

import (
	"context"
	"fmt"
)

// PaginatedResponse wraps a paginated API response.
type PaginatedResponse[T any] struct {
	Items       []T `json:"items"`
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
	CurrentPage int `json:"currentPage"`
}

// ListParams holds pagination and filtering parameters.
type ListParams struct {
	Page     int
	PageSize int
	Query    string
}

// DefaultListParams returns sensible defaults.
func DefaultListParams() ListParams {
	return ListParams{Page: 1, PageSize: 25}
}

// QueryParams converts ListParams to URL query parameters.
func (p ListParams) QueryParams() map[string]string {
	params := map[string]string{}
	if p.Page > 0 {
		params["page"] = fmt.Sprintf("%d", p.Page)
	}
	if p.PageSize > 0 {
		params["pageSize"] = fmt.Sprintf("%d", p.PageSize)
	}
	if p.Query != "" {
		params["query"] = p.Query
	}
	return params
}

// AutoPaginate fetches all pages using the provided function.
func AutoPaginate[T any](ctx context.Context, fn func(params ListParams) (*PaginatedResponse[T], error), base ListParams) ([]T, error) {
	if base.Page == 0 {
		base.Page = 1
	}
	if base.PageSize == 0 {
		base.PageSize = 100
	}

	var all []T
	for {
		resp, err := fn(base)
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Items...)
		if base.Page >= resp.TotalPages {
			break
		}
		base.Page++
	}
	return all, nil
}
