# Implementation

## Summary
- Implemented:
  - Added `scripts/trigger-error-spike.sh`, a small incident replay script that appends valid ERROR JSON Lines to `./logs/demo-node/app.log` by default.
  - Updated `POST /api/alerts/config` to accept partial positive config updates, so `{"threshold":5}` works for quick demos while omitted fields keep their current values.
  - Preserved validation for explicitly invalid config values and added rejection for empty `{}` requests.
  - Added handler tests for threshold-only config and empty config rejection.
  - Fixed review finding: partial config requests now pass only explicitly provided positive fields to `AlertEngine.UpdateConfig`, avoiding handler-level read-modify-write with a stale config snapshot.
  - Added a regression test proving sequential partial updates preserve earlier fields.
  - Updated README and docs to describe partial config updates and the incident replay flow.
- Changed files:
  - `api-server/handlers/alerts.go`
  - `api-server/handlers/alerts_test.go`
  - `scripts/trigger-error-spike.sh`
  - `README.md`
  - `docs/api.md`
  - `docs/project-roadmap.md`
  - `docs/testing-evidence.md`
  - `docs/report-notes.md`
  - `docs/deployment.md`
  - `docs/one-month-defense-roadmap.md`
  - `docs/knowledge-base.md`

## Deviations From Plan
- Deviation: Docker runtime end-to-end verification happened after the initial implementation pass, during follow-up verification.
- Reason: Static validation/script dry-run came first; runtime testing then caught and fixed the host-file permission issue in the trigger script.

## Notes For Verification
- Behavior to verify:
  - `POST /api/alerts/config` with `{"threshold":5}` returns `200` and preserves `window_seconds`/`cooldown_seconds`.
  - `POST /api/alerts/config` with `{}` returns `400`.
  - `./scripts/trigger-error-spike.sh 20` writes 20 JSON Lines with `level=ERROR` and `service=demo-node` to `./logs/demo-node/app.log`.
  - With Docker stack running, Filebeat ships the scripted ERROR lines, Elasticsearch count increases, and alerting emits `[alerting] alert sent`.
- Known limitations:
  - The script intentionally supports simple ASCII-safe messages only; custom `MESSAGE` values containing quotes, backslashes, or newlines are rejected to avoid malformed JSON.
  - Runtime dashboard/WebSocket alert verification remains pending until Docker services are running.

## Next Handoff
- Current phase: implementation
- Next phase: review
- Must read: `agent-context/features/mvp-completion-scope/04-verification.md`, `scripts/trigger-error-spike.sh`, `api-server/handlers/alerts.go`, `docs/testing-evidence.md`
- Decisions locked: Incident replay is implemented as a local script, not a new demo-only API endpoint; alert config endpoint supports partial positive updates.
- Open risks: Dashboard screenshots are still useful for final slides/report, but core runtime alert evidence is captured.
- Validation status: Static checks and Docker runtime E2E passed after the review finding fix and script fallback fix.
