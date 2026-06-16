# Discovery

## Repo Facts
- Demo services emit JSON Lines with fields `timestamp`, `level`, `service`, `message`, and `metadata`.
- Filebeat ships each JSON line to Logstash in the Filebeat `message` field.
- Logstash parses `message` into `[log]`, then should promote fields to root for the Go API.
- The Go API filters `level.keyword` and reads root fields into `LogEntry`.

## Current Failure
- Dashboard level cells are empty because `/api/logs` returns `level: ""` and `log_message: ""`.
- Elasticsearch has documents, but a sample document keeps `level` and `message` under `[log]` and contains `_mutate_error` plus `_dateparsefailure`.
- Logstash logs show: `wrong argument type String (expected LogStash::Timestamp)`.

## Root Cause
- `logstash/pipeline/logstash.conf` renames `[log][timestamp]` directly into `@timestamp`.
- `@timestamp` expects a Logstash timestamp object, not the original JSON string.
- When that mutate block fails, the remaining root field promotion is not reliable, so API/dashboard cannot see `level` and `message`.

## Constraints
- Keep JSON Lines as the canonical demo log format.
- Keep Go API response shape unchanged.
- Keep the fix scoped to Logstash field promotion and timestamp parsing.

## Next Handoff
- Current phase: discovery
- Next phase: implementation
- Must read: `agent-context/features/compose-logstash-healthcheck/02-plan.md`, `logstash/pipeline/logstash.conf`, `api-server/handlers/logs.go`
- Decisions locked: Fix the pipeline, not the dashboard.
- Open risks: Old bad ES documents will remain unless explicitly reindexed/deleted.
- Validation status: Failure reproduced with ES/API probes.
