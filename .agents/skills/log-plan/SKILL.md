---
name: log-plan
description: Plan a log-system feature before implementation. Use when the user asks to analyze, scope, design, or create a decision-complete plan for changes in this repository, including API, alerting, pipeline, dashboard, demo services, Docker Compose, or Elasticsearch behavior.
---

# Log Plan

Load shared project rules from `.agents/skills/log-system-dev/SKILL.md`, then use this skill as the planning entrypoint.

Read:

- `.agents/GUIDE.md`
- `.agents/skills/log-system-dev/references/phase-discovery.md`
- The subsystem phase reference that matches the requested feature
- The applicable `AGENTS.md` files

## Required Context Output

For every non-trivial feature, create or update `.agents/context/features/<feature-slug>/` unless the user explicitly says not to persist context.

If `.agents/context` is read-only, create or update `agent-context/features/<feature-slug>/` instead. Use the same phase filenames and treat this as persisted context, not as `not persisted`.

Use a short lowercase hyphen slug based on the feature, for example `api-logs-count` or `dashboard-alert-banner`.

If both `.agents/context` and `agent-context` are unavailable, do not claim context was written. Return the exact intended paths and content, mark the context as `not persisted`, and make the next writable phase create those files.

Write:

- `01-discovery.md`: repo facts, applicable `AGENTS.md` constraints, current implementation state, relevant files, unknowns, and risks.
- `02-plan.md`: chosen implementation plan, accepted API/UI/data contracts, acceptance criteria, assumptions, and `Next Handoff`.

Use templates from `.agents/context/templates/` when readable; otherwise follow the same headings from existing files under `agent-context`. Do not change application source files during planning unless the user explicitly asks for implementation.

Return a concise implementation plan with summary, key changes, test plan, assumptions, and the context folder path written.
