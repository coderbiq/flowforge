# Workflow Rules

Use one primary mode per turn.

## Start

1. Resolve the task with `flowforge task ready --type i` or an explicit task id.
2. Confirm the task is `ready` and is an implementation task.
3. Run `flowforge context task --task <id>`.
4. Read only via CLI. Do not open wiki files or `.flowforge` files directly.
5. Confirm linked requirement, design, and constraints are present in task context.

## Refuse Early

Stop immediately if the task is `not_ready`, `blocked`, `design`, or `analysis`.
Report the blocking status and do not start code changes.

## Execute

1. Make only the changes needed for the ready implementation task.
2. Re-check the task context through CLI before writing if the scope is unclear.
3. Record important events with `flowforge log create --kind progress --title <title> --for <card-id> --summary <text>`.
4. Do not edit card files, frontmatter, wikilinks, generated navigation, or internal card links manually.
5. If a task needs extra knowledge, use `library suggest` / `card search --scope library` and link confirmed context through CLI.
6. Hand-written Markdown links are allowed only for external source references.

## Finish

Run project tests and `flowforge validate all` when card state changed.

Report:

- changed files
- tests run or not run
- validation result if card state changed
- gaps or blockers
- one next step
