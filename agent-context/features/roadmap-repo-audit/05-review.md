# Review

## Findings
- Severity: Medium
- File/line: `api-server/handlers/alerts.go:71`, `api-server/alerting/engine.go:172`, `api-server/alerting/engine.go:175`
- Issue: `POST /api/alerts/config` validates only `threshold`. A request with an invalid `window_seconds` or `cooldown_seconds` still returns `200 {"status":"updated"}` and can partially update other fields. Verified with `{"threshold":5,"window_seconds":-1,"cooldown_seconds":60}` returning 200 while keeping the old window value.
- Recommendation: Validate all provided config fields in the handler before calling `UpdateConfig`, and return `400` for non-positive `window_seconds` or `cooldown_seconds` instead of silently ignoring them.
- Status: Fixed in `api-server/handlers/alerts.go` and covered by `api-server/handlers/alerts_test.go`.

- Severity: Low
- File/line: `docker-compose.yml:82`, `docker-compose.yml:101`, `docker-compose.yml:119`
- Issue: `filebeat`, `demo-node`, and `demo-go` do not define `healthcheck`, even though the project Compose rules require every service to expose one. `docker compose ps` shows these services as only `Up`, while `elasticsearch`, `logstash`, and `api-server` report healthy.
- Recommendation: Add lightweight healthchecks: Filebeat can check its process or HTTP endpoint if enabled; demo services can check that their log file exists and is writable/non-empty.
- Status: Fixed in `docker-compose.yml`; runtime `docker compose ps` now reports all six services healthy.

## Test Gaps
- Gap: There are still no automated tests for Elasticsearch query construction.
- Risk: Runtime smoke tests pass today, but a future field-contract regression could still require Docker smoke testing to catch.

- Gap: Runtime verification currently requires Docker socket access and local ports outside the sandbox.
- Risk: Review can verify this manually, but the repo does not yet encode the end-to-end checks as repeatable automated tests.

## Residual Risks
- Risk: Dashboard browser console reports `404 /favicon.ico`. This is cosmetic and does not block dashboard operation, but it leaves avoidable console noise.
- Risk: CORS allows `*` and WebSocket origin checks always return true. This remains acceptable for local development, but should become env-configurable before production exposure.

## Commit Readiness
- Ready: Yes
- Reason: Both review findings are fixed and verified. Remaining risks are cosmetic or production-hardening follow-ups outside this fix.

## Verification
- `docker compose config` passed.
- `GOCACHE=/tmp/log-system-go-build-cache go test ./...` from `api-server` passed.
- `GOCACHE=/tmp/log-system-demo-go-build-cache go test ./...` from `services/demo-go` passed.
- `node --check services/demo-node/index.js` passed.
- `node --check dashboard/app.js` passed.
- `docker compose up -d` started all services; `docker compose ps` showed Elasticsearch, Logstash, and API healthy.
- After fixes, `docker compose up -d --build` recreated containers and `docker compose ps` showed Elasticsearch, Logstash, API, Filebeat, demo-node, and demo-go healthy.
- Invalid `POST /api/alerts/config` with `window_seconds:-1` now returns `400`.
- Valid alert config still returns `200`, and defaults were restored after the runtime test.
- `GET /api/health` returned 200 with Elasticsearch connected.
- `GET /api/logs?size=3` returned non-empty `@timestamp`, `level`, `service`, and `log_message`.
- `GET /api/logs/count` returned nonzero INFO/WARN/ERROR counts.
- Elasticsearch `logs-*` count returned nonzero documents.
- Playwright loaded the dashboard, confirmed WebSocket connected, rendered stats/table/pagination, and reported only the favicon 404 console error.

## Next Handoff
- Current phase: review
- Next phase: log-git
- Must read: `agent-context/features/roadmap-repo-audit/03-implementation.md`, `agent-context/features/roadmap-repo-audit/04-verification.md`, `agent-context/features/roadmap-repo-audit/05-review.md`
- Decisions locked: API/log pipeline field contract is currently working in runtime smoke tests.
- Open risks: favicon route is absent; production CORS/WebSocket origin hardening remains future work.
- Validation status: Static checks and runtime smoke tests passed after fixes.
