# Verification

## Commands
- Command: `docker compose build --progress=plain`
- Result: Pass
- Important output: `api-server`, `demo-go`, and `demo-node` images built successfully.

- Command: `docker compose up -d`
- Result: Pass after fix
- Important output: `log-logstash` changed from unhealthy to healthy, and the full stack started.

- Command: `docker compose ps`
- Result: Pass
- Important output: `log-logstash` is `Up ... (healthy)`; `log-elasticsearch` and `log-api-server` are also healthy.

## Cause
- The Logstash healthcheck in `docker-compose.yml` was probing `http://localhost:9600/_node/stats`.
- Inside the container, that HTTP probe returned an empty reply and Docker kept the service unhealthy.
- The runtime itself was fine; Logstash had already started its API endpoint and was accepting the TCP socket on port 9600.

## Fix
- Replaced the Logstash healthcheck with a direct TCP socket probe against `127.0.0.1:9600`.
- This avoids the false unhealthy state caused by the HTTP probe.

## Validation
- `docker compose up -d` completed successfully.
- `docker compose ps` shows `log-logstash` as healthy and the stack running.

## Pipeline Field Contract Follow-up
- Command: `docker compose config`
- Result: Pass
- Important output: Compose config renders successfully after Logstash pipeline and healthcheck changes.

- Command: `docker compose restart logstash`
- Result: Pass
- Important output: `log-logstash` restarted and returned to healthy state.

- Command: `docker logs --since 1m log-logstash`
- Result: Pass
- Important output: New startup logs show the pipeline running; no new `wrong argument type String (expected LogStash::Timestamp)` messages after the final pipeline restart.

- Command: `curl -sS 'http://localhost:8080/api/logs?size=5'`
- Result: Pass
- Important output: API now returns rows with populated `level` and `log_message`.

- Command: `curl -sS 'http://localhost:8080/api/logs/count'`
- Result: Pass
- Important output: API count returned non-zero level totals, e.g. `ERROR`, `INFO`, and `WARN`.

- Command: `curl -sS 'http://localhost:8080/api/logs?level=ERROR&size=3'`
- Result: Pass
- Important output: API returned ERROR rows with populated `level`, `service`, `log_message`, and `metadata`.

- Command: `curl -sS -u elastic:changeme123 'http://localhost:9200/logs-*/_search?size=1&sort=@timestamp:desc'`
- Result: Pass
- Important output: New ES document has root `level`, `service`, `log_message`, `metadata`, canonical `@timestamp`, and no `_mutate_error` / `_dateparsefailure` tag.
