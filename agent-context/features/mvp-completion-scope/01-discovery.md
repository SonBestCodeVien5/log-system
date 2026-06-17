# Discovery

## Request
- Feature/request: Identify what remains to complete the MVP besides adding an active alert trigger script.
- Feature slug: mvp-completion-scope

## Repo Facts
- Current state: Core MVP code is present for Docker Compose, demo services, Filebeat, Logstash, Elasticsearch, Go API, alerting engine, and dashboard.
- Current state: Project docs already frame the next phase as end-to-end verification, evidence capture, defense preparation, and a small incident replay script.
- Current state: `docs/testing-evidence.md` still has all Step 10 checks marked pending, so the MVP is not yet closed by evidence.
- Relevant files: `docs/project-roadmap.md`, `docs/testing-evidence.md`, `docs/deployment.md`, `docs/one-month-defense-roadmap.md`, `api-server/handlers/alerts.go`, `api-server/handlers/alerts_test.go`, `api-server/alerting/engine.go`, `dashboard/app.js`, `docker-compose.yml`, `services/demo-node/index.js`, `services/demo-go/main.go`, `filebeat/filebeat.yml`, `logstash/pipeline/logstash.conf`.
- Existing constraints: Keep JSON Lines format unchanged, keep dashboard vanilla HTML/CSS/JS, do not add broad new scope before verification, do not hardcode secrets, use environment variables for runtime config.
- Existing constraints: Root `AGENTS.md` documents `AlertEngine.shouldAlert` must use one lock for atomic check/write; current implementation uses one `sync.Mutex` for config/dedup and a separate client mutex.
- Existing constraints: Worktree already has unrelated dashboard/screenshot changes; planning should not overwrite them.

## Applicable Instructions
- Root `AGENTS.md`: MVP path is demo services -> Filebeat -> Logstash -> Elasticsearch -> Go API -> dashboard/WebSocket alerts; response format is `{"data":[...],"total":...,"page":...,"size":...}`; dashboard pagination is 20 rows.
- Area `AGENTS.md`: API must use gin/go-elasticsearch and JSON errors; dashboard must support table, filters, alert banner, auto-refresh, threshold control; services must emit JSON Lines with uppercase level and stable fields.
- Skill references: `log-plan`, `phase-discovery`, `phase-api`, `phase-alerting`, `phase-pipeline`, `phase-dashboard`, `phase-services`, `phase-verification`.

## Unknowns And Risks
- Unknowns: Runtime status was not verified in this planning turn; Docker stack may still expose environment-specific issues.
- Risks: `docs/project-roadmap.md` suggests `POST /api/alerts/config` with only `{"threshold":5}`, but `api-server/handlers/alerts.go` currently rejects omitted `window_seconds` and `cooldown_seconds`; demo instructions can fail unless the command uses the full config or the handler is changed to support partial updates.
- Risks: Without an incident replay script, alert demo depends on random ERROR rate and can be flaky.
- Risks: Without real outputs, screenshots, and response-time numbers, docs remain descriptive rather than defensible.

## Next Handoff
- Current phase: discovery
- Next phase: plan
- Must read: `agent-context/features/mvp-completion-scope/02-plan.md`, `docs/project-roadmap.md`, `docs/testing-evidence.md`, `docs/one-month-defense-roadmap.md`.
- Decisions locked: Do not expand MVP scope beyond incident replay, contract cleanup, verification evidence, docs, screenshots, and clean-clone test.
- Open risks: Alert config contract/docs mismatch; runtime checks still pending.
- Validation status: Planning only; no runtime validation performed.
