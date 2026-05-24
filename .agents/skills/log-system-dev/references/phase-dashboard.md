# Phase: Dashboard

Use this context for `dashboard/index.html`, `dashboard/app.js`, and `dashboard/style.css`.

## Constraints

- No React, Vue, Tailwind, Bootstrap, or frontend build step.
- Keep logic in `app.js`, structure in `index.html`, and styling in `style.css`.
- Use `fetch()` for REST calls.
- Use `WebSocket` for alerts.

## Expected Features

- Log table with 20 rows per page by default.
- Filter bar: level, service/app, time/search where relevant.
- Alert banner that can be dismissed.
- Auto-refresh every 10 seconds.
- Threshold control that calls `POST /api/alerts/config`.

## Visual Rules

- Make log levels visually distinct:
  - INFO: `#2563eb`
  - WARN: `#d97706`
  - ERROR: `#dc2626`
- Keep the UI operational and dense enough for repeated log inspection.
- Avoid marketing-style landing pages; the dashboard itself is the first screen.
