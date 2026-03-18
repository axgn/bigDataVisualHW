package pbdbapi

import (
	"context"
	"net/url"
)

// IntervalsService provides access to /intervals endpoints.
type IntervalsService struct {
	client *Client
}

// IntervalsParams defines query parameters for interval list queries.
type IntervalsParams struct {
	Name   string
	Type   string
	Show   []string
	Limit  int
	Offset int
}

func (p IntervalsParams) values() url.Values {
	v := make(url.Values)
	setString(v, "name", p.Name)
	setString(v, "type", p.Type)
	setCSV(v, "show", p.Show)
	applyPagination(v, pagination{Limit: p.Limit, Offset: p.Offset})
	return v
}

// List fetches a list of intervals.
func (s *IntervalsService) List(ctx context.Context, p IntervalsParams) (*IntervalsResponse, error) {
	var out IntervalsResponse
	if err := s.client.doJSON(ctx, "/intervals/list.json", p.values(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
