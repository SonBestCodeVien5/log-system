# Log System Agent Guide

This folder stores Codex skills and working context for the `log-system` repo.

## Folder Roles

- `skills/`: Codex skills that appear in skill discovery and can be invoked with `$skill-name`.
- `context/`: feature-specific and phase-specific handoff notes created while working.
- `context/features/`: one folder per active feature, for example `features/api-logs-count/`.
- `context/phases/`: reusable phase notes that are not tied to one feature.
- `context/archive/`: completed feature contexts after they are no longer active.

## Skill Roles

- `$log-system-dev`: core project skill with shared rules, architecture constraints, and phase references.
- `$log-plan`: planning entrypoint.
- `$log-implement`: implementation entrypoint.
- `$log-debug`: debugging entrypoint.
- `$log-review`: review entrypoint.
- `$log-git`: post-review Git handoff and commit entrypoint.
- `$log-docs`: report-ready documentation and evidence entrypoint.

The wrapper skills intentionally reuse `$log-system-dev` as the shared baseline so project rules are defined once. They are not duplicates; they choose the workflow for the current intent, then load the relevant phase reference.

## Recommended Flow

1. Start with `$log-plan`.
2. `$log-plan` must create `.agents/context/features/<feature-slug>/` unless the user explicitly says not to persist context.
3. `$log-plan` must write `01-discovery.md` with repo facts, constraints, unknowns, and relevant files.
4. `$log-plan` must write `02-plan.md` with the chosen implementation plan and acceptance criteria.
5. Use `$log-implement` and record important implementation notes in `03-implementation.md`.
6. Run validation and record results in `04-verification.md`.
7. Use `$log-review` and record findings or residual risks in `05-review.md`.
8. Use `$log-git` after review to inspect status, prepare staging/commit, and record `06-git.md`.
9. Use `$log-docs` before large features to prepare docs, and after each feature to update report notes, decisions, and testing evidence while context is fresh.

Each phase file should end with:

```text
## Next Handoff
- Current phase:
- Next phase:
- Must read:
- Decisions locked:
- Open risks:
- Validation status:
```

## Feature Context Files

For every non-trivial feature, create:

```text
.agents/context/features/<feature-slug>/
```

Use these files inside it:

```text
01-discovery.md
02-plan.md
03-implementation.md
04-verification.md
05-review.md
06-git.md
07-blocked.md
```

Use `.agents/context/templates/` as the field guide for each file. Copy the template headings and fill only useful, task-specific content.

## Phase Responsibilities

- `$log-plan` creates the feature folder and writes `01-discovery.md` plus `02-plan.md`.
- `$log-implement` reads `01-discovery.md` and `02-plan.md` before editing, then writes `03-implementation.md`.
- Verification work writes `04-verification.md` after checks, even when checks fail.
- `$log-review` reads prior phase files, then writes `05-review.md`.
- `$log-git` reads review context, then writes `06-git.md`.
- `$log-docs` reads feature context and `docs/*.md`, then updates report-ready documentation and evidence.
- `$log-debug` writes `07-blocked.md` when a phase cannot complete or a runtime failure blocks progress.

For trivial one-line fixes, context may be skipped only when the user explicitly says not to persist context or the change has no meaningful follow-up phase.

If the active collaboration mode or permissions forbid file writes, do not claim context was persisted. Instead, include the exact intended context file paths and content in the response, mark the context as `not persisted`, and write it as soon as a later phase permits file changes.

## If A Phase Breaks

Create or update the feature's `07-blocked.md`.

Record:

- failed phase
- exact blocker
- evidence from files or commands
- attempted fixes
- rollback status
- needed user input or permission
- safest next action

Then continue with `$log-debug` for technical blockers or `$log-plan` if the scope needs to change.
