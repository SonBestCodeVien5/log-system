---
name: log-debug
description: Debug log-system failures. Use when logs do not reach Elasticsearch, Docker Compose services fail, API queries fail, WebSocket alerts do not arrive, or the dashboard does not show expected log or alert data.
---

# Log Debug

Load shared project rules from `.agents/skills/log-system-dev/SKILL.md`, then use this skill as the debugging entrypoint.

Read:

- `.agents/GUIDE.md`
- `.agents/skills/log-system-dev/references/phase-verification.md`
- The relevant subsystem phase reference
- `docs/deployment.md` for runtime pipeline issues

Follow the evidence path before proposing fixes: service log file, Filebeat harvester, Logstash input and Grok parse, Elasticsearch index, API query, dashboard fetch/WebSocket.

## Required Context Output

When a phase fails, a runtime path is blocked, or a fix cannot be completed, locate or create `.agents/context/features/<feature-slug>/` and write or update `07-blocked.md`.

If `.agents/context` is read-only, use `agent-context/features/<feature-slug>/07-blocked.md` instead. Treat this as persisted context.

If neither `.agents/context` nor `agent-context` is writable, return the intended `07-blocked.md` content and mark it as `not persisted`.

Record:

- failed phase or failing runtime path
- exact blocker
- evidence from files or commands
- attempted fixes
- rollback status
- required user input, permission, or environment change
- safest next action

If debugging succeeds, update `04-verification.md` or `03-implementation.md` with the resolved cause and validation result instead of leaving stale blocker context.

Use templates from `.agents/context/templates/` when readable; otherwise follow the same headings from existing files under `agent-context`.
