package pbdbapi

import (
	"context"
	"net/url"
)

// OccurrencesService provides access to /occs endpoints.
type OccurrencesService struct {
	client *Client
}

// OccurrencesParams defines query parameters for occurrence list queries.
type OccurrencesParams struct {
	TaxonID  int
	Interval string
	Location string
	Show     []string
	Limit    int
	Offset   int
}

func (p OccurrencesParams) values() url.Values {
	v := make(url.Values)
	setInt(v, "taxon_id", p.TaxonID)
	setString(v, "interval", p.Interval)
	setString(v, "loc", p.Location)
	setCSV(v, "show", p.Show)
	applyPagination(v, pagination{Limit: p.Limit, Offset: p.Offset})
	return v
}

// List fetches a list of occurrences.
func (s *OccurrencesService) List(ctx context.Context, p OccurrencesParams) (*OccurrencesResponse, error) {
	var out OccurrencesResponse
	if err := s.client.doJSON(ctx, "/occs/list.json", p.values(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
