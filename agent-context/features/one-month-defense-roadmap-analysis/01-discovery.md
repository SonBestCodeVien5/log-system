# Discovery

## Request
- Feature/request: Analyze `docs/one-month-defense-roadmap.md` and provide conclusions.
- Feature slug: `one-month-defense-roadmap-analysis`

## Repo Facts
- Current state:
  - `docs/one-month-defense-roadmap.md` is a defense-readiness plan, not a code implementation plan.
  - The roadmap prioritizes stable MVP, evidence, code understanding, a small personal incident replay feature, and final rehearsal.
  - Current repo already has Docker Compose, demo services, Filebeat, Logstash, Elasticsearch, Go API, alerting, and dashboard implemented.
  - Recent review/fix context says API/log pipeline and dashboard smoke tests passed, and the remaining review findings were fixed in the current uncommitted work.
  - Elasticsearch security is enabled in Compose, so ES verification commands should use `-u elastic:$ES_PASSWORD`.
  - `POST /api/alerts/config` now expects complete positive config values, so roadmap examples that send only `{"threshold":5}` are stale.
- Relevant files:
  - `docs/one-month-defense-roadmap.md`
  - `docs/architecture.md`
  - `docs/api.md`
  - `docs/deployment.md`
  - `docs/testing-evidence.md`
  - `docs/report-notes.md`
  - `docs/decisions.md`
  - `docker-compose.yml`
  - `api-server/handlers/alerts.go`
- Existing constraints:
  - Keep MVP architecture intact.
  - Do not add broad scope in the final month.
  - Preserve JSON Lines demo log format.
  - Dashboard remains vanilla HTML/CSS/JS.

## Applicable Instructions
- Root `AGENTS.md`:
  - Use JSON Lines as canonical demo log format.
  - Runtime config should come from env.
  - Dashboard should use vanilla HTML/JS/CSS.
  - Feature context belongs under `.agents/context` or `agent-context`, not `docs/`.
- Area `AGENTS.md`:
  - Root docs are project/report docs; no app source edits requested for this analysis.
- Skill references:
  - `$log-plan` requires discovery and plan context for non-trivial analysis.
  - `phase-discovery.md` says to inspect docs/architecture, docs/api, and docs/deployment for broad planning.

## Roadmap Observations
- Strengths:
  - The roadmap has the right strategic priority: evidence and understanding before new features.
  - It gives a concrete 4-week progression: verify MVP, study implementation, add a small incident replay, then rehearse/packaging.
  - It explicitly rejects risky late scope such as RBAC, Kubernetes, tracing, AI, dashboard rewrite, and database changes.
  - The recommended incident replay direction A is low-risk and aligned with the pipeline contract.
- Drift or inconsistencies:
  - ES curl examples are unauthenticated, but current Compose enables Elasticsearch security.
  - Alert config examples send only `threshold`, while the current API contract sends `threshold`, `window_seconds`, and `cooldown_seconds`.
  - The roadmap says to write answers for 8 defense questions, but the actual list contains 10 questions.
  - Some code line anchors in the code map are likely stale after implementation changes.
  - The roadmap treats several items as future verification work even though recent runtime smoke tests already confirmed the core MVP path.

## Unknowns And Risks
- Unknowns:
  - Whether `docs/testing-evidence.md`, `docs/report-notes.md`, and `docs/decisions.md` currently contain enough fresh evidence for defense.
  - Whether incident replay already exists elsewhere outside the reviewed roadmap file.
- Risks:
  - If roadmap commands stay stale, defense rehearsal may fail on ES auth or alert config validation.
  - If evidence is not captured from fresh runtime runs, the project may appear theoretical despite the stack working.
  - Adding the optional incident replay too late could create avoidable instability unless kept as a script-only slice.

## Next Handoff
- Current phase: discovery
- Next phase: plan
- Must read: `docs/one-month-defense-roadmap.md`, `docs/api.md`, `docs/deployment.md`, `agent-context/features/one-month-defense-roadmap-analysis/01-discovery.md`
- Decisions locked: analyze only; do not edit roadmap content in this phase.
- Open risks: stale commands and missing/old evidence need a docs update pass later.
- Validation status: document analysis only; no runtime validation run in this turn.
