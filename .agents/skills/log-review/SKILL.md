---
name: log-review
description: Review log-system changes. Use when the user asks for a code review, risk analysis, regression check, or pre-commit review for this repository, especially across Go API, alerting concurrency, pipeline config, dashboard/API contract, or demo service log format.
---

# Log Review

Load shared project rules from `.agents/skills/log-system-dev/SKILL.md`, then use this skill as the review entrypoint.

Read:

- `.agents/GUIDE.md`
- `.agents/skills/log-system-dev/references/phase-verification.md`
- Relevant subsystem phase references
- Applicable `AGENTS.md` files

Lead with findings ordered by severity. Focus on broken contracts, hardcoded config, incorrect log format, unhandled errors, race-prone shared state, missing validation, and test gaps.

## Required Context Input And Output

Before reviewing a non-trivial feature, read the active feature context under `.agents/context/features/<feature-slug>/`, especially:

- `01-discovery.md`
- `02-plan.md`
- `03-implementation.md`
- `04-verification.md`

After review, write or update `05-review.md` with:

- findings ordered by severity
- fixes applied or recommended
- test gaps
- residual risks
- commit readiness
- `Next Handoff`

Set `Next Handoff` to `log-git` when changes are ready for staging or commit.

If review is blocked, write or update `07-blocked.md`.

If writes are not permitted, return the intended `05-review.md` or `07-blocked.md` content and mark it as `not persisted`.

Use templates from `.agents/context/templates/`.
