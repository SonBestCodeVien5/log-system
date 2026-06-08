# Skill Usage Flow

Use this guide when starting, continuing, debugging, or recovering a `log-system` feature.

## How To Invoke Skills

Use the small wrapper skill that matches the current intent:

- `$log-plan`: analyze and plan a feature before editing.
- `$log-implement`: implement a scoped feature or fix.
- `$log-debug`: investigate a failing runtime, pipeline, API, alerting, or dashboard path.
- `$log-review`: review changes for defects and test gaps.
- `$log-git`: inspect Git status, prepare staging/commit, and write post-review Git handoff.
- `$log-docs`: update report-ready docs, decisions, and testing evidence from feature context.
- `$log-system-dev`: load the general project context directly.

If the UI exposes slash entries for these skills, they map to the same wrappers. If not, use the `$skill-name` form.

## Where Context Lives

Use two different locations:

- `.agents/GUIDE.md`: human-facing guide for folder roles, skill roles, and phase flow.
- `.agents/skills/`: stable instructions that should apply to every future session.
- `.agents/context/`: temporary or feature-specific working context that should carry one feature across phases.
- `agent-context/`: writable mirror for feature context when `.agents/context` is read-only in a Codex session.

Do not put feature-specific investigation notes into `SKILL.md`; that makes the skill stale and noisy.

## Feature Context Folder

For every non-trivial feature, create:

```text
.agents/context/features/<feature-slug>/
```

If `.agents/context` is read-only, create the same folder under:

```text
agent-context/features/<feature-slug>/
```

Inside it, use these phase files:

```text
01-discovery.md
02-plan.md
03-implementation.md
04-verification.md
05-review.md
06-git.md
07-blocked.md
```

Use short slugs such as `api-logs-count`, `dashboard-alert-banner`, or `pipeline-grok-parse`.

This context creation is mandatory by default for `$log-plan`. Skip it only when the user explicitly says not to persist context.

If `.agents/context` is read-only but `agent-context` is writable, write to `agent-context` and treat the handoff as persisted. If both locations are unavailable, do not claim the context was persisted. Return the exact intended file paths and content, mark the handoff as `not persisted`, and write it when file changes are allowed.

## Phase Flow

1. Discovery: record current repo state, relevant files, constraints from `AGENTS.md`, and unknowns.
2. Plan: record final decisions, scope, API/UI contracts, and acceptance criteria.
3. Implementation: record changed files, important implementation notes, and deviations from plan.
4. Verification: record commands run, outputs that matter, failures, and skipped tests.
5. Review: record findings, fixes applied, remaining risks, and commit readiness.
6. Git: record status, staging intent, commit message or hash, and remaining working tree state.
7. Docs: update report notes, decisions, testing evidence, and public docs while feature context is fresh.

Use `.agents/context/templates/` for the exact headings of each phase file when readable. If unavailable, mirror the same headings from existing files under `agent-context`.

At the end of each phase, append a `Next Handoff` section with:

```text
## Next Handoff
- Current phase:
- Next phase:
- Must read:
- Decisions locked:
- Open risks:
- Validation status:
```

The next phase starts by reading the previous phase file and its `Next Handoff`.

## What To Write For The Next Phase

Write only context that changes future decisions:

- file paths and symbols that matter
- accepted API or UI contracts
- constraints and non-goals
- failed approaches and why they failed
- commands already run and their result
- blockers that require user input or environment changes

Avoid copying long command output. Summarize the important lines and keep exact commands.

## If A Phase Fails

When a phase cannot complete, write or update `07-blocked.md` in the feature context folder.

Include:

- the failed phase
- exact blocker
- evidence from files or commands
- attempted fixes
- rollback status
- what input or permission is needed
- safest next action

Then switch to `$log-debug` if the blocker is technical, or `$log-plan` if the blocker changes scope.

## Phase Context Parking

The folders under `.agents/context/phases/` are parking areas for reusable phase notes that are not tied to one feature. Prefer `.agents/context/features/<feature-slug>/` for normal feature work, or `agent-context/features/<feature-slug>/` when `.agents/context` is read-only.
