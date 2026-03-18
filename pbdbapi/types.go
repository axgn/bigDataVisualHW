package pbdbapi

// ListResponse is a generic PBDB list response shape.
type ListResponse[T any] struct {
	Records      []T    `json:"records"`
	RecordsFound int    `json:"records_found,omitempty"`
	RecordsShown int    `json:"records_shown,omitempty"`
	Offset       int    `json:"offset,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	Warnings     any    `json:"warnings,omitempty"`
	DataURL      string `json:"data_url,omitempty"`
}

type CollectionRecord map[string]any
type TaxonRecord map[string]any
type OccurrenceRecord map[string]any
type IntervalRecord map[string]any

type CollectionsResponse = ListResponse[CollectionRecord]
type TaxaResponse = ListResponse[TaxonRecord]
type OccurrencesResponse = ListResponse[OccurrenceRecord]
type IntervalsResponse = ListResponse[IntervalRecord]
