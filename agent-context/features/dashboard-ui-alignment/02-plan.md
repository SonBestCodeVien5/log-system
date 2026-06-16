# Plan

## Implementation Steps
- Simplify the alert config strip into a single white tool bar with the same border/radius/padding language as the rest of the dashboard.
- Replace the custom grid layout with inline fields and one primary action button.
- Reduce visual weight on the config title and status text so they read as supporting UI, not a separate module.
- Keep the current `updateAlertConfig()` behavior and WebSocket sync flow unchanged.

## Acceptance Criteria
- The alert config section no longer looks like a separate panel.
- The alert config controls visually match the rest of the dashboard controls.
- Desktop layout remains clean and stable.
- The dashboard still loads and shows logs without layout overlap.

## Validation Commands
- `node --check dashboard/app.js`
- `docker compose up -d`
- Playwright snapshot of the dashboard on desktop

## Next Handoff
- Current phase: plan
- Next phase: implementation
- Must read: `dashboard/index.html`, `dashboard/app.js`, `dashboard/style.css`
- Decisions locked: Keep the current data and API behavior.
- Open risks: Visual tuning may still need a small follow-up after the first pass.
- Validation status: Pending.
