# Review

## Findings
- Severity: Medium
- File/line: `api-server/handlers/alerts.go:82-105`
- Issue: Partial alert config updates are merged in the HTTP handler by calling `GetConfig()`, modifying the returned copy, then calling `UpdateConfig()`. Because `GetConfig` and `UpdateConfig` take separate locks, two concurrent partial requests can lose each other's changes. Example: request A sends `{"threshold":5}` while request B sends `{"window_seconds":60}`; both can read the old config, then the later full update can overwrite the earlier field back to the old value.
- Recommendation: Avoid read-modify-write in the handler. Build an `alerting.AlertConfig` containing only explicitly provided positive fields and pass it to `AlertEngine.UpdateConfig`, which already updates only positive fields under one lock. Then call `GetConfig()` only after the update to build the response. Add a concurrent partial-update test or at least a test that sequential partial updates preserve prior changes.

## Test Gaps
- Gap: Runtime Docker E2E test was previously missing; now covered for script -> Filebeat/Logstash/Elasticsearch -> API query -> API alert log.
- Risk: Dashboard banner visual proof still needs a screenshot/browser check if required for report evidence.
- Gap: Multiple partial config updates are now covered by `TestUpdateConfigPreservesPriorPartialUpdates`.
- Risk: Concurrent request interleavings are not directly stress-tested, but the handler no longer performs read-modify-write and delegates partial updates to the engine lock.

## Residual Risks
- Risk: `scripts/trigger-error-spike.sh` intentionally rejects custom messages with quotes, backslashes, or newlines to keep JSON generation simple. This is acceptable for the default MVP demo, but should be documented as a script limitation if custom messages become part of the workflow.
- Risk: `docs/testing-evidence.md` now has runtime command evidence; screenshots remain optional report polish.

## Commit Readiness
- Ready: Yes, after the follow-up implementation fix.
- Reason: The partial config update lost-update issue was fixed by removing handler-level read-modify-write and adding regression coverage for preserving earlier partial updates. Runtime E2E alert evidence has also been captured through API/log checks.

## Next Handoff
- Current phase: review
- Next phase: log-git
- Must read: `api-server/handlers/alerts.go`, `api-server/alerting/engine.go`, `api-server/handlers/alerts_test.go`, `agent-context/features/mvp-completion-scope/03-implementation.md`, `agent-context/features/mvp-completion-scope/04-verification.md`
- Decisions locked: Incident replay script approach remains appropriate; API should keep partial positive update support.
- Open risks: Capture dashboard screenshots before final defense if needed.
- Validation status: Review finding fixed; static checks and runtime E2E alert path passed; ready for `log-git` if user wants to stage/commit this slice.
