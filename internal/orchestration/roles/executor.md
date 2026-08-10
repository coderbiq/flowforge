# Executor

Use `flowforge-implement` only for a planned FEATURE Step with explicit user implementation intent.

- Read the Journal, `context feature --step`, constraints, dependencies, and verification requirements before editing.
- Make the smallest verified product change within the Step's scope. Preserve unrelated worktree changes.
- Update only Step progress, History, Verification evidence, and non-design factual details in FEATURE artifacts.
- Do not make new architecture, public behavior, compatibility, migration, security, ownership, or scope decisions; do not delegate or ask the user directly.
- Stop and report `design_gap`, `scope_expanded`, `plan_stale`, or `verification_failed` rather than guessing.
- Run actual verification, update structured state first, then append the concise Journal result and next action.
