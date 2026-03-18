package pbdbapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTaxaParamsValues(t *testing.T) {
	p := TaxaParams{
		Name:     "Canis",
		Rank:     "genus",
		Interval: "Pleistocene",
		Show:     []string{"attr", "size"},
		Limit:    25,
		Offset:   50,
	}

	v := p.values()
	if got := v.Get("name"); got != "Canis" {
		t.Fatalf("name = %q, want %q", got, "Canis")
	}
	if got := v.Get("rank"); got != "genus" {
		t.Fatalf("rank = %q, want %q", got, "genus")
	}
	if got := v.Get("limit"); got != "25" {
		t.Fatalf("limit = %q, want %q", got, "25")
	}
	if got := v.Get("offset"); got != "50" {
		t.Fatalf("offset = %q, want %q", got, "50")
	}
	if got := v.Get("show"); got != "attr,size" {
		t.Fatalf("show = %q, want %q", got, "attr,size")
	}
}

func TestDoJSONRetriesOn500(t *testing.T) {
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"temporary"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"records":[{"id":1}]}`))
	}))
	defer ts.Close()

	c, err := NewClient(
		WithBaseURL(ts.URL),
		WithRetry(3, time.Millisecond),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, err := c.Collections().List(context.Background(), CollectionsParams{})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(resp.Records) != 1 {
		t.Fatalf("records = %d, want 1", len(resp.Records))
	}
	if attempts != 3 {
		t.Fatalf("attempts = %d, want 3", attempts)
	}
}

func TestDoJSONNoRetryOn400(t *testing.T) {
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer ts.Close()

	c, err := NewClient(
		WithBaseURL(ts.URL),
		WithRetry(3, time.Millisecond),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = c.Taxa().List(context.Background(), TaxaParams{Name: "x"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var se *StatusError
	if !errors.As(err, &se) {
		t.Fatalf("error type = %T, want *StatusError", err)
	}
	if se.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", se.StatusCode, http.StatusBadRequest)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}
