---
name: log-git
description: Prepare and perform Git handoff after log-system review. Use after implementation, verification, and review when the user asks to inspect status, prepare staging, create a commit, summarize changes, or produce a Git-ready handoff for this repository.
---

# Log Git

Load shared project rules from `.agents/skills/log-system-dev/SKILL.md`, then use this skill as the post-review Git entrypoint.

Read:

- `.agents/GUIDE.md`
- Active feature context under `.agents/context/features/<feature-slug>/` or fallback `agent-context/features/<feature-slug>/`
- `03-implementation.md`
- `04-verification.md`
- `05-review.md`

## Required Git Checks

Before any Git action, inspect:

- `git status --short`
- `git diff --stat`
- focused `git diff` for files relevant to the feature

Preserve unrelated user changes. Do not stage or commit files that are outside the feature scope unless the user explicitly includes them.

## Staging And Commit Rules

- Do not run `git add` or `git commit` unless the user explicitly asks to stage or commit.
- If committing, stage only reviewed feature files and any required context files.
- Use a concise conventional commit message when the repo has no stricter convention.
- Mention any uncommitted unrelated files that remain outside the commit.

## Required Context Input And Output

Before Git handoff, read the active feature context, especially `05-review.md`.

Write or update `06-git.md` with:

- git status summary
- files intended for staging
- files intentionally excluded
- commit message proposal or actual commit hash
- remaining working tree state
- `Next Handoff`

If Git handoff is blocked, write or update `07-blocked.md`.

If `.agents/context` is read-only, write to `agent-context/features/<feature-slug>/` instead and treat it as persisted context. If neither location is writable, return the intended `06-git.md` or `07-blocked.md` content and mark it as `not persisted`.

Use templates from `.agents/context/templates/` when readable; otherwise follow the same headings from existing files under `agent-context`.
