# Phase: Verification, Review, And Debugging

Use this context after implementation or when the user asks for review/debugging.

## Static Checks

- `docker compose config`
- `GOCACHE=/tmp/log-system-go-build-cache go test ./...` inside Go modules
- `node --check services/demo-node/index.js`

## Runtime Checks

Use these when Docker services are expected to run:

- `docker compose ps`
- `curl http://localhost:9200/_cluster/health`
- `curl "http://localhost:9200/logs-*/_count"`
- `curl http://localhost:8080/api/health`
- API query examples from `docs/api.md`

## Debug Order

For missing logs, check:

1. Demo service is writing `./logs/<service>/app.log`.
2. Filebeat is harvesting the file.
3. Logstash is receiving Beats events.
4. Grok is not producing `_grokparsefailure`.
5. Elasticsearch has `logs-*` indices.
6. API queries the same index and fields.

## Review Focus

Lead with defects and risks:

- Broken pipeline contracts
- Hardcoded config that should be env-driven
- Incorrect log format
- Missing error handling
- Race-prone alerting state
- Dashboard/API contract mismatches
