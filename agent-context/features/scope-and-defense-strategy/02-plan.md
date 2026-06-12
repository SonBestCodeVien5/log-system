# Plan

## Summary
- Goal: Choose a defensible direction for the final phase of the log-system project.
- Success criteria:
  - Student can explain every critical system path without relying on AI.
  - Repo has real end-to-end evidence, not only planned tests.
  - Any new scope is small, observable, and directly strengthens the defense story.

## Key Changes
- Implementation:
  - Primary recommendation: spend most remaining effort on deep understanding, runtime verification, docs consistency, evidence, and defense practice.
  - Optional scope expansion only after the current MVP is verified: add one small "personal value" feature that strengthens operations value without changing architecture.
  - Best optional feature candidates:
    - Incident replay / demo scenario: deterministic script or mode that generates a clear error spike for demo and testing.
    - Alert explanation panel or docs-backed runbook: show why alert fired, threshold, count, window, and suggested operator action.
    - Retention/ILM documentation or lightweight config: explain production log lifecycle without overbuilding.
    - Observability evidence pack: screenshots, command outputs, response times, and architecture decision records.
- Public API/UI/data contracts:
  - No API/UI contract changes are required for the recommended path.
  - If adding a small feature later, keep existing JSON Lines, `logs-*`, REST response shape, WebSocket alert contract, and vanilla dashboard constraints.
- Out of scope:
  - Major new systems such as auth, multi-tenant RBAC, Kubernetes deployment, AI log analysis, distributed tracing, or production-grade security unless the deadline leaves enough time and the student can defend them.
  - Broad refactors or stack changes.

## Acceptance Criteria
- Scenario: Student is asked why Go, Elasticsearch, Filebeat/Logstash, sliding window, dedup, JSON Lines, and dashboard pagination were chosen.
- Expected result: Student can answer with tradeoffs, code references, and measured evidence.
- Scenario: Project is demoed from clean/local environment.
- Expected result: Dashboard loads, logs flow from services to Elasticsearch, API filters work, alert threshold can trigger a WebSocket banner, and evidence is recorded.
- Scenario: A small personal feature is added after verification.
- Expected result: It has a narrow demo path, tests or manual evidence, and a defense explanation that fits within 1-2 minutes.

## Assumptions
- Assumption: This is a one-month graduation project where reliability of defense matters more than raw feature count.
- Assumption: Current MVP is close enough that polishing and evidence will produce more grading value than broad expansion.
- Assumption: The user wants honest strategic advice, not immediate code implementation.

## One-Month Direction
- Week 1: Freeze the MVP behavior and verify end-to-end. Run the Docker Compose path, confirm logs reach Elasticsearch, confirm API filters/counts, trigger WebSocket alerting, record real command output and response times.
- Week 2: Deep-study the implementation and repair documentation drift. Build a personal explanation map for Filebeat, Logstash, Elasticsearch Query DSL, Go API handlers, alerting mutex/cooldown, WebSocket, and dashboard pagination.
- Week 3: Add exactly one small personal-value feature if Week 1 verification is stable. Preferred choice: incident replay/demo scenario plus alert explanation/runbook because it improves demo reliability and is easy to defend.
- Week 4: Package for defense. Prepare slides/report screenshots, clean README flow, test clone from scratch, rehearse a 5-minute demo and a 10-minute Q&A.

## Recommended Personal Feature
- Chosen direction: Incident replay plus alert explanation.
- Why this over bigger scope: It shows operational thinking, makes the defense demo deterministic, and stays within the existing architecture.
- Possible implementation shape:
  - Add a controlled demo scenario that generates a burst of ERROR logs without changing the normal JSON Lines contract.
  - Surface or document alert context: count, threshold, window, cooldown, and suggested operator action.
  - Record evidence in `docs/testing-evidence.md` and explain the scenario in `docs/report-notes.md`.
- Non-goal: Do not add auth, Kubernetes, tracing, AI analysis, or a new database in the final month unless all verification and defense materials are already strong.

## Next Handoff
- Current phase: plan
- Next phase: verification/docs or optional feature planning
- Must read: `docs/project-roadmap.md`, `docs/testing-evidence.md`, `docs/report-notes.md`, `docs/knowledge-base.md`
- Decisions locked: Recommended allocation is approximately 50% deep understanding and defense evidence, 25% verification/docs polish, 15% one small personal feature, 10% rehearsal and cleanup.
- Open risks: Timeline and grading rubric may change the exact allocation.
- Validation status: Planning context persisted; no source files changed.
