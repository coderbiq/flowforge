# Workflow Rules

Use one primary mode per turn.

## Start

1. Resolve the task with `flowforge task ready --type i` or an explicit task id.
2. Confirm the task is `ready` and is an implementation task.
3. Run `flowforge context task --task <id>`.
4. Read only via CLI. Do not open wiki files or `.flowforge` files directly.

## Refuse Early

Stop immediately if the task is `not_ready`, `blocked`, `design`, or `analysis`.
Report the blocking status and do not start code changes.

## Execute

1. Make only the changes needed for the ready implementation task.
2. Re-check the task context through CLI before writing if the scope is unclear.
3. Record important events with `flowforge log create --kind progress --title <title> --for <card-id> --summary <text>`.

## Finish

Report:

- changed files
- tests run or not run
- gaps or blockers
- one next step
