# Project Guidelines

## Code Style
- Use idiomatic Go and keep package-level APIs small and explicit.
- Prefer the existing functional options style for client configuration (`Option func(*Client) error`).
- Keep service methods thin: build query params and delegate HTTP/JSON handling to client helpers.
- Preserve current error model (`StatusError`, `APIError`, `DecodeError`) and wrap errors with `%w` when propagating.

## Architecture
- This repository is a lightweight PBDB API client centered on `Client` in `client.go`.
- Endpoint groupings live behind service accessors in `services.go`:
  - `CollectionsService` -> `/colls/list.json`
  - `TaxaService` -> `/taxa/list.json`
  - `OccurrencesService` -> `/occs/list.json`
  - `IntervalsService` -> `/intervals/list.json`
- Query construction should stay in per-service `*Params.values()` methods and shared helpers in `query.go`.
- Response payloads intentionally use flexible record types (`map[string]any` aliases in `types.go`) because PBDB fields vary by request.

## Build and Test
- Run all tests: `go test ./...`
- Build all packages: `go build ./...`
- Keep module metadata clean: `go mod tidy`
- When changing HTTP behavior (retry, timeout, status handling), add or update `httptest`-based coverage in `client_test.go`.

## Conventions
- Default client behavior uses:
  - Base URL `https://paleobiodb.org/data1.2`
  - Timeout `30s`
  - Retries enabled for `429` and `5xx` only
- Do not add client-side validation for all PBDB query semantics unless requested; current design defers most parameter validation to the API.
- Keep retry behavior deterministic and context-aware (`sleepWithContext` + exponential backoff).
- Maintain backward compatibility for exported symbols unless a breaking change is explicitly requested.

## Pitfalls
- Do not retry `4xx` responses by default.
- Preserve `ctx` cancellation behavior in request and backoff paths.
- `WithTimeout` only applies when `httpClient` is `*http.Client`; avoid assuming custom clients honor timeout options.
