---
name: log-implement
description: Implement a log-system feature using the repo conventions. Use when the user asks Codex to build or modify API handlers, alerting, Filebeat or Logstash pipeline, Elasticsearch config, Docker Compose, demo services, or the vanilla dashboard in this repository.
---

# Log Implement

Load shared project rules from `.agents/skills/log-system-dev/SKILL.md`, then use this skill as the implementation entrypoint.

Read the subsystem phase reference that matches the task:

- Workflow guide: `.agents/GUIDE.md`
- API: `.agents/skills/log-system-dev/references/phase-api.md`
- Alerting: `.agents/skills/log-system-dev/references/phase-alerting.md`
- Pipeline/Docker: `.agents/skills/log-system-dev/references/phase-pipeline.md`
- Dashboard: `.agents/skills/log-system-dev/references/phase-dashboard.md`
- Demo services: `.agents/skills/log-system-dev/references/phase-services.md`

## Required Context Input And Output

Before editing, locate the active feature folder under `.agents/context/features/<feature-slug>/` and read:

- `01-discovery.md`
- `02-plan.md`

If no feature context exists for non-trivial work, create it first using `$log-plan` behavior or write a minimal `01-discovery.md` and `02-plan.md` before implementation.

If prior plan context was returned as `not persisted`, create the missing context files before editing application code.

After implementation, write or update:

- `03-implementation.md`: changed files, behavior implemented, deviations from plan, unresolved risks, and `Next Handoff`.
- `04-verification.md`: validation commands and results when checks are run during implementation.

Use templates from `.agents/context/templates/`.

Inspect current files first, preserve unrelated user changes, implement the smallest complete vertical slice, then validate with `.agents/skills/log-system-dev/references/phase-verification.md`.
