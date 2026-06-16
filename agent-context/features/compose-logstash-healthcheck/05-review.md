# Review

## Findings
- High: Resolved. Logstash no longer renames the JSON string timestamp directly into `@timestamp`; it parses through a temporary field and promotes `level`, `service`, `log_message`, and `metadata` to the API/dashboard contract.
- Medium: Resolved. `services/AGENTS.md`, agent phase references, and report-facing docs now describe JSON Lines plus `log_message` consistently.
- Low: Existing Elasticsearch documents indexed before the pipeline fix can still contain `_mutate_error` / `_dateparsefailure` and missing root fields. New documents are correct.

## Fixes Applied Or Recommended
- Applied: `docker-compose.yml` Logstash healthcheck now uses a TCP socket probe against `127.0.0.1:9600`, avoiding the false unhealthy state from the previous HTTP probe.
- Applied: `logstash/pipeline/logstash.conf` copies `[log][timestamp]` to a temporary `log_timestamp`, parses it into `@timestamp`, removes the temporary field, and promotes JSON message text to `log_message`.
- Applied: updated stale service-specific agent/docs text that still described the old bracket log format.

## Test Gaps
- `docker compose config`: pass.
- `node --check services/demo-node/index.js`: pass.
- `GOCACHE=/tmp/log-system-go-build-cache go test ./...` in `api-server`: pass, but packages currently have no test files.
- `docker compose ps`: pass; API, Elasticsearch, and Logstash are healthy.
- API/ES runtime probes: `/api/health`, `/api/logs/count`, `/api/logs?level=ERROR&size=3`, and latest ES document are healthy and match the dashboard contract.

## Residual Risks
- Old ES documents remain in the index with the previous bad shape. For clean demo evidence, delete or reindex old `logs-*` data deliberately.
- Documentation now matches the JSON Lines pipeline; no known doc drift remains for this field contract.

## Commit Readiness
- Ready for staging/commit for the healthcheck plus pipeline contract fix.

## Next Handoff
- Current phase: review
- Next phase: log-git
- Must read: `agent-context/features/compose-logstash-healthcheck/04-verification.md`, `agent-context/features/compose-logstash-healthcheck/05-review.md`, `docker-compose.yml`, `logstash/pipeline/logstash.conf`
- Decisions locked: Keep JSON Lines as canonical log format.
- Open risks: Old ES documents may still show the old bad shape until cleaned or ignored.
- Validation status: healthcheck and pipeline/API field contract validated.
