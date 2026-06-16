# Implementation

## Changed Files
- `docker-compose.yml`
- `logstash/pipeline/logstash.conf`

## Behavior Implemented
- Logstash healthcheck now uses a TCP probe against `127.0.0.1:9600`, avoiding the false unhealthy result from the previous HTTP probe.
- Logstash no longer renames the JSON string timestamp directly to `@timestamp`.
- The JSON timestamp is copied to `log_timestamp`, parsed by the `date` filter into canonical Elasticsearch `@timestamp`, then the temporary field is removed.
- JSON log fields are promoted to the API/dashboard contract:
  - `[log][level]` -> `level`
  - `[log][service]` -> `service`
  - `[log][message]` -> `log_message`
  - `[log][metadata]` -> `metadata`
- Grok enrich now reads from `log_message`.

## Deviations From Plan
- The first implementation promoted `[log][message]` to `message`, which fixed `level` but left dashboard messages blank because the API response contract uses `log_message`.
- The final implementation promotes to `log_message` so no API rebuild is required.

## Unresolved Risks
- Old Elasticsearch documents indexed before the fix still contain `_mutate_error` / `_dateparsefailure` and may have missing root fields.
- New documents are correct. For a clean demo dataset, old indices should be deleted or reindexed deliberately.

## Next Handoff
- Current phase: implementation
- Next phase: verification
- Must read: `agent-context/features/compose-logstash-healthcheck/04-verification.md`, `logstash/pipeline/logstash.conf`
- Decisions locked: Keep `@timestamp` as the canonical parsed timestamp field; do not keep the temporary `log_timestamp`.
- Open risks: Existing bad ES documents remain.
- Validation status: Implemented and runtime-verified.
