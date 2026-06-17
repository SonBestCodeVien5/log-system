# Plan

## Summary
- Goal: turn the roadmap analysis into a decision-ready set of conclusions and follow-up actions.
- Success criteria: identify whether the roadmap is useful, where it is stale, what should be prioritized next, and what should not be added before defense.

## Key Changes
- Implementation: no code or docs implementation in this planning turn.
- Public API/UI/data contracts:
  - Preserve current API contracts from `docs/api.md`.
  - Treat complete alert config body as the accepted contract: `threshold`, `window_seconds`, `cooldown_seconds`.
  - Treat authenticated ES curl as the accepted deployment verification shape when security is enabled.
- Out of scope:
  - Editing `docs/one-month-defense-roadmap.md`.
  - Implementing incident replay.
  - Running a full clean-clone test.

## Conclusions To Present
- The roadmap is directionally strong and should remain the main final-month operating plan.
- The project should now treat week 1 mostly as evidence capture, not discovery or build work, because the MVP is already implemented and recently smoke-tested.
- The most valuable next feature remains script-based incident replay, not a new API endpoint or major product feature.
- The highest-priority docs update is to refresh commands and API examples so rehearsal uses the real current system.
- The roadmap should become a checklist with dated evidence links, screenshots, and measured numbers rather than only a plan.

## Acceptance Criteria
- Scenario: User asks whether the roadmap is still valid.
- Expected result: Answer distinguishes "strategically valid" from "some operational commands are stale".

- Scenario: User asks what to do next.
- Expected result: Recommend evidence capture, docs drift cleanup, then script-based incident replay.

- Scenario: User asks whether to add more features.
- Expected result: Recommend avoiding broad features and only adding controlled incident replay if evidence is already stable.

## Assumptions
- Assumption: Current uncommitted fixes for alert config validation and Compose healthchecks are intended to stay.
- Assumption: Defense priority is reliability and explainability over feature breadth.
- Assumption: The user wants analysis and conclusions, not immediate document edits.

## Next Handoff
- Current phase: plan
- Next phase: log-docs
- Must read: `agent-context/features/one-month-defense-roadmap-analysis/01-discovery.md`, `agent-context/features/one-month-defense-roadmap-analysis/02-plan.md`, `docs/one-month-defense-roadmap.md`
- Decisions locked: keep incident replay script-first; avoid broad late-month features.
- Open risks: roadmap commands need update for ES auth and alert config full body.
- Validation status: context persisted; no source changes made.
