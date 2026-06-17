# Plan

## Summary
- Goal: Close the MVP by making alert testing repeatable, fixing or aligning the alert config demo contract, and recording end-to-end verification evidence.
- Success criteria: A fresh run can show logs flowing into Elasticsearch, API filters/counts working, dashboard rendering data, WebSocket alert appearing from a controlled incident, and docs containing real outputs/metrics/screenshots.

## Key Changes
- Implementation: Add a small incident replay script that appends 10-20 valid ERROR JSON Lines to the existing tailed log file, preferably `./logs/demo-node/app.log`.
- Implementation: Align `POST /api/alerts/config` contract with docs/demo. Either send full config everywhere or update handler/tests to allow partial positive updates while preserving validation for invalid explicit values.
- Implementation: Run and record Step 10 end-to-end checks: dashboard static serve, ES count, API health, API filters, level counts, alert trigger, response time.
- Implementation: Capture screenshots for normal dashboard, ERROR filter, and alert banner.
- Implementation: Run a clean-clone or clean-directory smoke test before final defense.
- Public API/UI/data contracts: Keep JSON Lines fields `timestamp`, `level`, `service`, `message`, `metadata`; keep ES/API contract `@timestamp`, `level`, `service`, `log_message`, `metadata`; keep dashboard/API endpoints unchanged unless only partial alert config support is added.
- Out of scope: Authentication/RBAC, Kubernetes, tracing, AI analysis, multi-tenant UI, dashboard framework migration, log format changes, Elasticsearch replacement, and large API refactors.

## Acceptance Criteria
- Scenario: Start stack with `docker compose up -d --build`.
- Expected result: `docker compose ps` shows core services healthy or running, and `curl http://localhost:8080/api/health` returns Elasticsearch connected.
- Scenario: Verify ingestion.
- Expected result: demo log files exist, Elasticsearch `logs-*` count is greater than zero, and `/api/logs?size=20` returns the standard paginated response.
- Scenario: Verify API filters.
- Expected result: `/api/logs?level=ERROR&size=3` returns only ERROR rows; `/api/logs?app=demo-node&size=3` returns only demo-node rows; `/api/logs/count` returns INFO/WARN/ERROR/total.
- Scenario: Verify dashboard.
- Expected result: `http://localhost:8080` serves the dashboard, table renders logs, filters update the table, stats update, and WebSocket status connects.
- Scenario: Verify alert.
- Expected result: Lower threshold using a valid config request, run incident replay script, then alert banner appears or API logs contain `[alerting] alert sent`; evidence includes command output and screenshot.
- Scenario: Verify docs.
- Expected result: `docs/testing-evidence.md`, `docs/deployment.md`, and `docs/report-notes.md` contain real measured results, not pending placeholders.

## Assumptions
- Assumption: The fastest safe incident replay approach is a local script appending valid JSON Lines to the existing log path rather than adding a demo API endpoint.
- Assumption: It is acceptable for final MVP to use local-development security defaults as long as `.env` remains uncommitted and docs tell users to change `ES_PASSWORD` before real deployment.
- Assumption: Runtime verification should happen after current unrelated dashboard changes are either accepted or intentionally handled.

## Next Handoff
- Current phase: plan
- Next phase: implementation
- Must read: `agent-context/features/mvp-completion-scope/01-discovery.md`, `docs/one-month-defense-roadmap.md` section "Incident Replay / Controlled Error Spike", `api-server/handlers/alerts.go`, `api-server/handlers/alerts_test.go`, `services/demo-node/index.js`, `filebeat/filebeat.yml`.
- Decisions locked: Use a small script for active trigger; keep MVP scope tight; prioritize evidence over new features.
- Open risks: Decide whether to fix the alert config endpoint to allow partial updates or update every demo command to send full config.
- Validation status: Not yet run; next phase should execute static/runtime checks and update evidence.
