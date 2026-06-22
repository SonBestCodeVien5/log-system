# Plan

## Summary
- Goal: Close the MVP as a defensible graduation deliverable by aligning docs with verified behavior, capturing the remaining evidence, and producing a report/slide package without expanding application scope.
- Success criteria: one clean-clone run is recorded; three report-ready screenshots exist; roadmap/architecture/deployment/report notes match current evidence; a coherent report outline and slide deck cover architecture, implementation, demo, measurements, limitations, and future work; demo and Q&A are rehearsed.

## Key Changes
- Implementation: freeze feature development. Optionally remove the favicon 404 only as cosmetic demo polish. Add focused tests for alert cooldown/dedup and logs handler query/validation only if schedule permits; treat these as confidence improvements, not MVP blockers.
- Evidence: run a clean-clone test in a separate directory and record setup time; capture normal dashboard, ERROR filter, and live alert banner screenshots; repeat incident replay after ensuring cooldown is clear and record actual alert latency.
- Documentation: update stale status tables and pending wording in `docs/project-roadmap.md`, `docs/architecture.md`, `docs/deployment.md`, and `docs/report-notes.md`; convert `report-notes.md` from placeholders into report-ready chapter material and link every claim to code/evidence.
- Slide/report package: prepare a 10-slide narrative: problem, requirements, architecture, log contract, ingest pipeline, REST/dashboard, alert algorithm, deterministic incident demo, measured results, limitations/future work. Reuse `docs/decisions.md`, `docs/testing-evidence.md`, and the screenshot pack as sources.
- Public API/UI/data contracts: no changes. Keep `/api/health`, `/api/logs`, `/api/logs/count`, `/api/alerts/config`, `/ws/alerts`, JSON Lines, `logs-*`, and the WebSocket alert shape stable.
- Out of scope: authentication/RBAC, ILM implementation, Kubernetes, distributed tracing, AI analysis, multi-tenancy, framework rewrite, database change, and broad refactors.

## Acceptance Criteria
- Scenario: A new checkout is started using README instructions.
- Expected result: all Compose services become healthy, dashboard loads with data, and elapsed setup time plus any environment prerequisite is recorded.
- Scenario: The five-minute demo is run.
- Expected result: filters work, incident replay creates searchable ERROR logs, WebSocket banner appears within an explained interval, and cooldown behavior does not surprise the presenter.
- Scenario: A reviewer reads roadmap, architecture, deployment, evidence, and report notes.
- Expected result: no document still claims completed runtime/evidence work is pending; measurements include date/environment and are not presented as load benchmarks.
- Scenario: Slides/report are reviewed against the code.
- Expected result: every feature claim has a code reference or evidence artifact; limitations clearly separate MVP/local-demo choices from production requirements.

## Delivery Order
1. P0: update stale docs status and create a single MVP closure checklist.
2. P0: perform and record clean-clone verification.
3. P0: capture three screenshots, especially the WebSocket alert banner.
4. P0: write the report-ready content and 10-slide deck from existing docs/evidence.
5. P1: rehearse the five-minute demo three times and prepare concise answers for the ten questions in the defense roadmap.
6. P2: add focused unit/race tests and favicon only if P0/P1 are complete.

## Test Plan
- Static: `docker compose config --quiet`; API and demo Go tests; Node and shell syntax checks; `git diff --check` after docs work.
- Runtime: `docker compose ps`; ES cluster/index count; `/api/health`; `/api/logs` level/service/search/pagination; `/api/logs/count`; dashboard console/network; partial alert config; incident replay; WebSocket banner.
- Reproducibility: execute the README flow from a separate fresh clone without relying on existing volumes or untracked files.

## Assumptions
- Assumption: the current repository behavior is the MVP contract and should be frozen.
- Assumption: screenshots and slide/report artifacts may be committed under a docs/evidence location once the required university format is known.
- Assumption: missing production auth and ILM are honest limitations/future work, not blockers for this local graduation MVP.

## Next Handoff
- Current phase: plan
- Next phase: `log-docs`, followed by focused verification
- Must read: `agent-context/features/final-mvp-report-audit/01-discovery.md`, this file, `docs/testing-evidence.md`, `docs/report-notes.md`, `docs/decisions.md`
- Decisions locked: freeze scope; prioritize truthful docs, reproducibility, screenshots, report/slide, and rehearsal.
- Open risks: university formatting/rubric is unknown; clean-clone and visual alert evidence are still unrecorded.
- Validation status: decision-complete plan persisted; no application source files changed.
