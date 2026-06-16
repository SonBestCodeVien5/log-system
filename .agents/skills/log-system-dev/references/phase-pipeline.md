# Phase: Pipeline, Elasticsearch, And Docker Compose

Use this context for `docker-compose.yml`, `filebeat`, `logstash`, and `elasticsearch` changes.

## Pipeline Contract

Data flow:

```text
demo services -> /logs/**/*.log -> Filebeat -> Logstash :5044 -> Elasticsearch :9200 -> Go API -> Dashboard
```

## Docker Compose Rules

- Every service uses `restart: always`.
- Every service has a healthcheck.
- All services join `log-network`.
- `./logs` is mounted for log producers and Filebeat.
- Filebeat and API wait for healthy upstream services where Compose supports it.

## Logstash Rules

- Input uses Beats on `${LOGSTASH_PORT}` or the configured container port.
- Logstash must parse JSON Lines emitted by demo services. The canonical input line is:

```json
{"timestamp":"2024-01-15T10:23:11Z","level":"ERROR","service":"demo-node","message":"Payment gateway timeout","metadata":{"order_id":"789"}}
```

- Promote parsed fields to the API/dashboard contract: `@timestamp`, `level`, `service`, `log_message`, and `metadata`.
- Grok is an optional enrich step over `log_message`; it must not drop logs that do not match.
- Output index should be date-based `logs-YYYY.MM.dd`.

## Elasticsearch Rules

- Version follows `ES_VERSION`.
- Port follows `ES_PORT`.
- Do not hardcode production credentials in committed files.
- Keep local dev security choices explicit in compose/config.

## Debug Path

Check in order: Elasticsearch health, Logstash logs, Filebeat harvester logs, actual files under `./logs`, Grok parse failures, and `logs-*` indices.
