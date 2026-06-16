---
name: log-system-dev
description: Develop, debug, review, and extend the log-system project. Use for tasks involving the Go API server, alerting engine, Elasticsearch queries, Filebeat/Logstash pipeline, demo log services, vanilla dashboard, Docker Compose, project feature planning, or phase-based implementation context for this repo.
---

# Log System Dev

## Core Workflow

For the human-facing folder and phase handoff guide, read `.agents/GUIDE.md`.

Start every task by reading the nearest `AGENTS.md` files that apply to the touched area:

- Root repo rules: `AGENTS.md`
- API work: `api-server/AGENTS.md`
- Dashboard work: `dashboard/AGENTS.md`
- Demo service work: `services/AGENTS.md`

Then load only the phase reference that matches the current task:

- Skill usage and context handoff flow: `references/usage-flow.md`
- Discovery and feature planning: `references/phase-discovery.md`
- API server and Elasticsearch handlers: `references/phase-api.md`
- Alerting and WebSocket behavior: `references/phase-alerting.md`
- Filebeat, Logstash, Elasticsearch, Docker Compose: `references/phase-pipeline.md`
- Dashboard UI and vanilla JS flows: `references/phase-dashboard.md`
- Demo log producers: `references/phase-services.md`
- Verification, review, and debugging: `references/phase-verification.md`

## Project Rules

Keep changes aligned with the repo's MVP architecture:

- Demo services write JSON Lines to `/logs/<service>/app.log`.
- Filebeat tails `/logs/**/*.log` and ships to Logstash.
- Logstash parses JSON Lines into Elasticsearch `logs-*`, then enriches `log_message` with Grok when patterns match.
- Go API exposes REST and WebSocket endpoints for the dashboard.
- Dashboard stays HTML, vanilla JS, and handwritten CSS only.

Use `.env` for runtime config. Do not hardcode passwords, ports, or hostnames in application code except static browser defaults already documented for local dashboard access.

## Implementation Defaults

- Prefer small, vertical feature slices: pipeline/config, API contract, dashboard behavior, then verification.
- Use `log.Printf` in Go server code; do not use `fmt.Println`.
- Return JSON errors as `{"error":"..."}` from HTTP handlers.
- Use `sync.RWMutex` for shared alerting state.
- Keep the demo service JSON Lines format unchanged because Logstash and the API/dashboard contract depend on it.
- Update `go.mod` when adding Go dependencies.

## Context Handoff

For feature work that spans multiple turns or phases, write handoff context under `.agents/context/features/<feature-slug>/` using the phase file names described in `.agents/GUIDE.md` and `references/usage-flow.md`.

If `.agents/context` is read-only in the current Codex session, use the writable mirror `agent-context/features/<feature-slug>/` with the same phase filenames instead. Keep stable skill instructions in `.agents/skills`; keep feature-specific discoveries, decisions, blockers, and verification notes in `.agents/context` or the `agent-context` fallback, never in `docs/`.
