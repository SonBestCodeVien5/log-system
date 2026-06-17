# Verification

## Commands
- Command: `docker compose config`
- Result: Passed.
- Important output: Compose rendered healthchecks for `filebeat`, `demo-node`, and `demo-go`.

- Command: `GOCACHE=/tmp/log-system-go-build-cache go test ./...` from `api-server`
- Result: Passed.
- Important output: `ok github.com/SonBestCodeVien5/log-system/api-server/handlers`.

- Command: `node --check services/demo-node/index.js`
- Result: Passed.
- Important output: no syntax errors.

- Command: `GOCACHE=/tmp/log-system-demo-go-build-cache go test ./...` from `services/demo-go`
- Result: Passed.
- Important output: package compiles with no test files.

- Command: `docker compose up -d --build`
- Result: Passed.
- Important output: rebuilt `api-server`, `demo-node`, and `demo-go`; recreated runtime containers.

- Command: `docker compose ps`
- Result: Passed.
- Important output: all services report healthy: `elasticsearch`, `logstash`, `api-server`, `filebeat`, `demo-node`, and `demo-go`.

## Scenarios
- Scenario: POST invalid alert config with `window_seconds:-1`.
- Expected: API returns `400` and does not accept partial update.
- Actual: `HTTP/1.1 400 Bad Request` with `{"error":"window_seconds must be >= 1"}`.

- Scenario: POST valid alert config.
- Expected: API returns `200` and current config matches request.
- Actual: `HTTP/1.1 200 OK` with threshold/window/cooldown from request.

- Scenario: Restore default alert config.
- Expected: API returns `200` with `threshold=10`, `window_seconds=300`, `cooldown_seconds=60`.
- Actual: `HTTP/1.1 200 OK` with default config values.

- Scenario: Query API logs after rebuild.
- Expected: log entries include non-empty `@timestamp`, `level`, `service`, and `log_message`.
- Actual: `GET /api/logs?size=3` returned populated entries.

- Scenario: Query API counts after rebuild.
- Expected: INFO/WARN/ERROR counts are nonzero when demo services are emitting.
- Actual: `GET /api/logs/count` returned nonzero counts and total.

## Failures Or Skips
- Failure/skip: Browser favicon 404 was not fixed.
- Reason: It is cosmetic and outside the two review findings requested for this implementation pass.

## Next Handoff
- Current phase: verification
- Next phase: review
- Must read: `agent-context/features/roadmap-repo-audit/03-implementation.md`, `agent-context/features/roadmap-repo-audit/04-verification.md`
- Decisions locked: The two prior review findings are fixed and verified.
- Open risks: Favicon route and production hardening for CORS/WebSocket origin remain future work.
- Validation status: Static checks, unit tests, Compose rebuild, healthchecks, and runtime API smoke tests passed.
