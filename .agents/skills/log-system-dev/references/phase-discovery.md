# Phase: Discovery And Feature Planning

Use this context before implementing a new feature or when the request is broad.

## Read First

- `AGENTS.md`
- The nearest area-specific `AGENTS.md`
- `docs/architecture.md`
- `docs/api.md` if the feature touches API or dashboard contracts
- `docs/deployment.md` if the feature touches Docker, ES, Filebeat, or Logstash

## Planning Checklist

- Identify the affected subsystem: pipeline, API, alerting, dashboard, demo services, or docs.
- Confirm the user-facing API or UI behavior before editing.
- Keep the MVP architecture intact unless the user explicitly asks to redesign it.
- Check for empty placeholder files before assuming existing implementation.
- Define the acceptance test as an observable path, for example: service writes log, Filebeat ships, ES indexes, API returns, dashboard renders.

## Output Shape

For plans, include summary, key changes, test plan, and assumptions. Mention only the files needed to avoid ambiguity.
