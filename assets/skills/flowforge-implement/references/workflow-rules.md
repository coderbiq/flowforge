# Workflow Rules

Use one primary mode per activation.

## Execute

1. Re-check the task context through CLI before writing if the scope is unclear.
2. Record important events with `flowforge log create --kind progress --title <title> --for <card-id> --summary <text>`.
3. If a task needs extra knowledge, use `library suggest` / `card search --scope library` and link confirmed context through CLI.
4. Hand-written Markdown links are allowed only for external source references.

## Finish

When the implementation is complete, report:

- changed files
- tests run
- validation result (if card state changed)
- gaps or blockers
- one next step
