package pbdbapi

import (
	"context"
	"net/url"
)

// TaxaService provides access to /taxa endpoints.
type TaxaService struct {
	client *Client
}

// TaxaParams defines query parameters for taxa list queries.
type TaxaParams struct {
	Name     string
	Rank     string
	Interval string
	Show     []string
	Limit    int
	Offset   int
}

func (p TaxaParams) values() url.Values {
	v := make(url.Values)
	setString(v, "name", p.Name)
	setString(v, "rank", p.Rank)
	setString(v, "interval", p.Interval)
	setCSV(v, "show", p.Show)
	applyPagination(v, pagination{Limit: p.Limit, Offset: p.Offset})
	return v
}

// List fetches a list of taxa.
func (s *TaxaService) List(ctx context.Context, p TaxaParams) (*TaxaResponse, error) {
	var out TaxaResponse
	if err := s.client.doJSON(ctx, "/taxa/list.json", p.values(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
