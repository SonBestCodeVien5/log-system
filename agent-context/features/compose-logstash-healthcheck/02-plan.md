# Plan

## Implementation Steps
- Change Logstash to copy `[log][timestamp]` into a temporary string field instead of renaming it to `@timestamp`.
- Promote `[log][level]`, `[log][service]`, `[log][message]`, and `[log][metadata]` to root fields.
- Parse the temporary timestamp with the `date` filter using `target => "@timestamp"`.
- Remove the temporary timestamp field after parsing.

## Acceptance Criteria
- Logstash stops emitting `wrong argument type String (expected LogStash::Timestamp)` for new events.
- New ES documents have root `level`, `service`, `message`, and `metadata`.
- `/api/logs?size=3` returns non-empty `level` and `log_message`.
- `/api/logs?level=ERROR&size=3` returns matching rows when new ERROR events arrive.
- `/api/logs/count` reports non-zero values once new promoted events are indexed.

## Validation Commands
- `docker compose config`
- `docker compose restart logstash`
- `docker compose ps`
- `docker logs --tail 80 log-logstash`
- `curl -sS 'http://localhost:8080/api/logs?size=3'`
- `curl -sS 'http://localhost:8080/api/logs/count'`

## Next Handoff
- Current phase: plan
- Next phase: implementation
- Must read: `agent-context/features/compose-logstash-healthcheck/01-discovery.md`, `logstash/pipeline/logstash.conf`
- Decisions locked: Keep the API/dashboard contract as-is.
- Open risks: Existing bad documents may still appear in broad queries until newer fixed events arrive.
- Validation status: Pending implementation.
