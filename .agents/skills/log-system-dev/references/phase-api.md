# Phase: Go API And Elasticsearch

Use this context for `api-server` changes.

## Required Conventions

- Module path: `github.com/SonBestCodeVien5/log-system/api-server`.
- Use gin for HTTP routing.
- Use go-elasticsearch v8 for Elasticsearch access.
- Read config from environment variables.
- Return errors as JSON and do not panic outside `main`.
- Response shape for list endpoints:

```json
{
  "data": [],
  "total": 0,
  "page": 1,
  "size": 20
}
```

## Endpoints

- `GET /api/health`
- `GET /api/logs`
- `GET /api/logs/count`
- `POST /api/alerts/config`
- `GET /ws/alerts`

## Elasticsearch Query Rules

- Query `logs-*`.
- Filter exact levels using `level.keyword`.
- Filter services using `service.keyword`.
- Use `@timestamp` range filters for time windows.
- Use `match` on `message` for full-text `q`.
- Sort logs by `@timestamp desc` unless the task says otherwise.

## Local Validation

Use `GOCACHE=/tmp/log-system-go-build-cache go test ./...` from `api-server` to avoid cache permission issues in sandboxed runs.
