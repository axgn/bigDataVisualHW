package pbdbapi

// Collections returns the service for collection endpoints.
func (c *Client) Collections() *CollectionsService {
	return &CollectionsService{client: c}
}

// Taxa returns the service for taxa endpoints.
func (c *Client) Taxa() *TaxaService {
	return &TaxaService{client: c}
}

// Occurrences returns the service for occurrence endpoints.
func (c *Client) Occurrences() *OccurrencesService {
	return &OccurrencesService{client: c}
}

// Intervals returns the service for interval endpoints.
func (c *Client) Intervals() *IntervalsService {
	return &IntervalsService{client: c}
}
