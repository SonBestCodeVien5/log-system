---
name: log-docs
description: Standardize log-system documentation for graduation reports and project handoff. Use when the user asks to update docs, summarize feature evidence, capture architecture/API/deployment changes, write report notes, record technical decisions, or prepare testing evidence from feature context.
---

# Log Docs

Load shared project rules from `.agents/skills/log-system-dev/SKILL.md`, then use this skill as the documentation and report-readiness entrypoint.

Read:

- `.agents/GUIDE.md`
- Active feature context under `.agents/context/features/<feature-slug>/` or fallback `agent-context/features/<feature-slug>/`
- `docs/architecture.md`
- `docs/api.md`
- `docs/deployment.md`
- `docs/report-notes.md`
- `docs/decisions.md`
- `docs/testing-evidence.md`
- Source files relevant to the documented feature

## Documentation Scope

Update docs when a feature changes or clarifies:

- system architecture, data flow, or component responsibility
- public API endpoints, query params, response shapes, or error behavior
- deployment, environment variables, Docker Compose, or pipeline setup
- verification commands, test results, runtime evidence, or known limitations
- technical decisions and rationale useful for a graduation report

Do not edit application code. If code and docs disagree, record the mismatch and recommend the code or docs change instead of silently inventing behavior.

## Required Outputs

Write or update the relevant docs:

- `docs/architecture.md`: architecture and data flow changes.
- `docs/api.md`: endpoint contracts and examples.
- `docs/deployment.md`: setup, operations, and troubleshooting changes.
- `docs/report-notes.md`: report-ready summary, feature narrative, screenshots/evidence placeholders, and explanation text.
- `docs/decisions.md`: technical decisions, alternatives considered, and rationale.
- `docs/testing-evidence.md`: commands, outcomes, evidence, and unresolved test gaps.

For feature-specific documentation work, also update `.agents/context/features/<feature-slug>/04-verification.md`, `05-review.md`, or `06-git.md` only when the docs work changes those handoffs. If `.agents/context` is read-only, update the same files under `agent-context/features/<feature-slug>/`.

If neither `.agents/context` nor `agent-context` is writable for feature context, return the intended file paths and content and mark them as `not persisted`.

## Writing Rules

- Keep docs factual and traceable to repo files, feature context, or command output.
- Prefer Vietnamese for report-facing explanations unless the surrounding file is already English-only.
- Preserve existing headings when possible and add sections only when they improve later report writing.
- Include dates only when they are part of captured evidence or commands, not as decoration.
- Do not overclaim completed behavior when the implementation is still scaffolded or unverified.
