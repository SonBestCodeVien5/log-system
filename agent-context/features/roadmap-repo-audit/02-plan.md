# Plan

## Summary
- Goal: bring roadmap and docs back in line with the actual repo state, after fixing the highest-priority pipeline/API field contract issue.
- Success criteria: API log list and count endpoints return meaningful `level`, `service`, and `message` data; roadmap marks current progress accurately; docs record real verification evidence.

## Key Changes
- Implementation: fix Logstash/API field contract before starting dashboard work.
- Public API/UI/data contracts: decide whether the canonical message field is `message` or `log_message`, then make Logstash, API structs, ES queries, docs, and dashboard plan agree.
- Out of scope: implementing dashboard UI in this planning phase.

## Proposed Roadmap State
- Done: bootstrap project, Docker Compose structure, demo services, Go compile/build verification, API server startup.
- Partial: Filebeat/Logstash pipeline, Go logs API, Docker Compose integration verification.
- Next: fix and verify Elasticsearch field contract across Logstash and API.
- Todo: dashboard log viewer, alert banner and threshold control, end-to-end alert test, documentation evidence, defense preparation.

## Acceptance Criteria
- Scenario: Query ES count.
- Expected result: `curl -u elastic:$ES_PASSWORD http://localhost:9200/logs-*/_count` returns count greater than 0.

- Scenario: Query API log count.
- Expected result: `curl http://localhost:8080/api/logs/count` returns nonzero counts for INFO/WARN/ERROR when logs exist in the time window.

- Scenario: Query API logs.
- Expected result: `curl "http://localhost:8080/api/logs?size=3"` returns entries with non-empty timestamp, level, service, and message fields.

- Scenario: Open dashboard root.
- Expected result: after dashboard implementation, `http://localhost:8080/` renders an operational log viewer rather than an empty file.

## Test Plan
- Run `GOCACHE=/tmp/log-system-go-build-cache go test ./...` from `api-server`.
- Run `docker compose build api-server`.
- Run `docker compose up -d`.
- Run `docker compose ps`.
- Run authenticated ES health and count checks.
- Run `/api/health`, `/api/logs/count`, and `/api/logs?size=3`.
- Inspect a raw ES document if API values are empty.

## Assumptions
- The project should keep JSON Lines as the primary log format.
- The dashboard should wait until API field contract is stable.
- The runtime password remains the local development default unless `.env` changes it.

## Next Handoff
- Current phase: plan
- Next phase: implementation
- Must read: `agent-context/features/roadmap-repo-audit/01-discovery.md`, `agent-context/features/roadmap-repo-audit/02-plan.md`, `logstash/pipeline/logstash.conf`, `api-server/handlers/logs.go`, `docs/api.md`, `actual_progress_roadmap.html`
- Decisions locked: do not start dashboard before fixing API field contract.
- Open risks: Logstash `_mutate_error` and `_dateparsefailure` root cause still needs implementation-phase debugging.
- Validation status: runtime stack is healthy, ES has logs, API health passes, API log fields/count are not yet correct.
