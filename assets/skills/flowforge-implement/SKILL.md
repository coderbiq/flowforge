---
name: flowforge-implement
description: Use ONLY when the user asks to execute a ready FlowForge implementation task, or provides a task id and wants code changes for that task. Do NOT use for design, analysis, feedback, archive, or general card lookup.
---

# flowforge-implement

## Start

Run `context feature --feature <id> --step <n>` to get the minimal execution context.
Confirm linked constraints and dependencies are present.

## Workflow

Token-aware execution loop:

1. Get step context: `context feature --feature <id> --step <n>` (only ~400 tokens)
2. Implement the step as described
3. Record: `card steps <id> --status done <n>` + `card log <id> --event "..."`
4. Validate: run tests + `flowforge validate all`
5. Next step or complete: `card evolve <id> --stage done` when all steps done

## Hard Rules

- Start each step with `context feature --feature <id> --step <n>`. This is your primary context.
- Never read the whole FEATURE card during step execution. Use section-level reading for supplemental info.
- Execute steps in order; skip blocked steps.
- After each step, use `card log` to record progress and `card steps` to update status.
- CLI for structured ops only (link, evolve, log, steps); direct file editing for body content.
- If implementation reveals a missing detail, edit the FEATURE card directly to add it.
- If a design issue is found: `card log --kind blocked` + `card steps --status blocked N --reason "..."`
- Run tests and `flowforge validate all` when card state changes.

## Output

Report changed files, tests run, validation result, gaps or blockers, and one next step.
