# Discovery

## Repo Facts
- Dashboard is vanilla HTML, JS, and CSS with a dense log table, stats bar, filter bar, WebSocket alert banner, and alert config inputs.
- Current alert config UI was moved into a separate strip, but it reads visually different from the rest of the dashboard.
- The dashboard already uses simple white cards, thin borders, small radii, and compact black action buttons.

## UI Problem
- The alert config area feels like a custom panel instead of part of the same system as the stats cards and filter bar.
- The grid layout, uppercase title, and inline status text make it stand out more than the rest of the page.

## Goal
- Make the alert config UI look like a normal dashboard tool strip.
- Keep it operational and aligned with the rest of the page rather than decorative.

## Constraints
- No framework or build step.
- Keep all logic in `dashboard/app.js`, structure in `dashboard/index.html`, styling in `dashboard/style.css`.
- Preserve the current API contract and WebSocket behavior.

## Next Handoff
- Current phase: discovery
- Next phase: implementation
- Must read: `dashboard/index.html`, `dashboard/app.js`, `dashboard/style.css`
- Decisions locked: Use the existing visual language, not a new theme.
- Open risks: The alert config status text could still feel too strong if styled like a badge instead of a small inline note.
- Validation status: Pending.
