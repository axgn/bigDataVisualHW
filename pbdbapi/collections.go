package pbdbapi

import (
	"context"
	"net/url"
)

// CollectionsService provides access to /colls endpoints.
type CollectionsService struct {
	client *Client
}

// CollectionsParams defines query parameters for collection list queries.
type CollectionsParams struct {
	ID       int
	BaseName string
	Show     []string
	Limit    int
	Offset   int
}

func (p CollectionsParams) values() url.Values {
	v := make(url.Values)
	setInt(v, "id", p.ID)
	setString(v, "base_name", p.BaseName)
	setCSV(v, "show", p.Show)
	applyPagination(v, pagination{Limit: p.Limit, Offset: p.Offset})
	return v
}

// List fetches a list of collections.
func (s *CollectionsService) List(ctx context.Context, p CollectionsParams) (*CollectionsResponse, error) {
	var out CollectionsResponse
	if err := s.client.doJSON(ctx, "/colls/list.json", p.values(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
