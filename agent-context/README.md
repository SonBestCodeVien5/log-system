# Agent Context

This folder stores feature-specific planning, implementation, verification,
review, and handoff notes when `.agents/context` is read-only in a Codex
session.

## Layout

Use the same phase filenames as `.agents/context`:

```text
agent-context/features/<feature-slug>/
├── 01-discovery.md
├── 02-plan.md
├── 03-implementation.md
├── 04-verification.md
├── 05-review.md
├── 06-git.md
└── 07-blocked.md
```

## Rule

Prefer `.agents/context/features/<feature-slug>/` when it is writable.
If `.agents/context` is read-only, persist context here instead of writing
feature context into `docs/`.
