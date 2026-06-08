# Discovery

## Request
- Feature/request: Scan the current repo, read project docs, and evaluate whether `actual_progress_roadmap.html` matches the real project state.
- Feature slug: `roadmap-repo-audit`

## Repo Facts
- The project is a centralized logging platform: demo services write JSON Lines logs, Filebeat ships to Logstash, Logstash indexes into Elasticsearch, Go API queries ES, and dashboard should render logs plus alerts.
- Runtime stack was observed running with `api-server`, `elasticsearch`, `logstash`, `filebeat`, `demo-node`, and `demo-go`.
- `GET /api/health` returned `{"elasticsearch":"connected","status":"ok"}`.
- Elasticsearch `logs-*` contained more than 11,000 documents during audit.
- `GET /api/logs?size=3` returned data and a nonzero total, but sampled entries had empty `level` and `log_message`.
- `GET /api/logs/count` returned zero counts for INFO, WARN, and ERROR despite ES containing logs.
- Raw ES documents showed fields under `log.level` and `log.message`, while API queries expect root fields such as `level.keyword`.
- Raw ES documents included `_mutate_error` and `_dateparsefailure`, indicating Logstash field promotion/timestamp parsing needs implementation-phase investigation.
- Dashboard files are placeholders: `dashboard/index.html`, `dashboard/app.js`, and `dashboard/style.css` are all 0 bytes.

## Relevant Files
- `actual_progress_roadmap.html`
- `docs/architecture.md`
- `docs/api.md`
- `docs/deployment.md`
- `docs/knowledge-base.md`
- `docs/testing-evidence.md`
- `docs/report-notes.md`
- `docs/decisions.md`
- `docker-compose.yml`
- `logstash/pipeline/logstash.conf`
- `filebeat/filebeat.yml`
- `api-server/main.go`
- `api-server/handlers/logs.go`
- `api-server/handlers/alerts.go`
- `api-server/alerting/engine.go`
- `dashboard/index.html`
- `dashboard/app.js`
- `dashboard/style.css`

## Applicable Instructions
- Root `AGENTS.md`: keep JSON Lines log format, use env config, API response shape should be `{"data":[...],"total":N,"page":N,"size":N}`, dashboard remains vanilla HTML/CSS/JS.
- Dashboard `AGENTS.md`: dashboard must include log table, filters, pagination, alert banner, auto-refresh, and threshold control.
- API `AGENTS.md`: API owns Gin routes, ES client, logs handlers, alert handlers, and alerting engine.
- `$log-plan`: planning only; write discovery and plan context when writable.

## Roadmap Assessment
- The roadmap is directionally aligned with the project goals.
- The roadmap is stale for current status.
- Step 6, `go mod tidy + verify compile`, should no longer be marked as next because Go tests/build and Docker build have already passed.
- Step 7, `docker compose up + verify pipeline`, should be marked partial because the stack runs and ES has logs, but API field contract is broken.
- Steps 8 and 9 are correctly todo because dashboard files are empty.
- The summary counts are inconsistent: the header says 13 total steps, but the visible roadmap lists 12 steps.
- Step 11 references `report.tex`, but no `report.tex` was found in the repo.

## Unknowns And Risks
- Unknown: whether Logstash `mutate rename` fails because `@timestamp` cannot be renamed from nested field or because the Filebeat event shape differs from the expected shape.
- Risk: dashboard implementation will be blocked or misleading until `/api/logs` and `/api/logs/count` return correct fields.
- Risk: docs currently include stale operational commands, such as unauthenticated Elasticsearch curl examples while ES security is enabled.
- Risk: `docs/knowledge-base.md` includes older explanations that conflict with current JSON Lines and alert locking conventions.

## Next Handoff
- Current phase: discovery
- Next phase: plan
- Must read: `agent-context/features/roadmap-repo-audit/01-discovery.md`, `actual_progress_roadmap.html`, `logstash/pipeline/logstash.conf`, `api-server/handlers/logs.go`
- Decisions locked: dashboard should wait until API field contract is stable.
- Open risks: Logstash `_mutate_error` and `_dateparsefailure` root cause still needs debugging.
- Validation status: runtime stack is healthy, ES has logs, API health passes, API log fields/count are not yet correct.
