# Workflow Rules

This reference defines the v3 feedback loop. Keep discovery, classification,
routing, and verification separate; use the CLI for every structured write.

## Quick Reference

| Command | Purpose |
|---|---|
| `flowforge card init --type feature --title "..." --proposal <id>` | Create an implementation scope |
| `flowforge card log <id> --kind bug --event "..."` | Record reproducible evidence |
| `flowforge journal append --actor flowforge-feedback --message "..." --references <id>` | Record proposal feedback |
| `flowforge card steps <id> --status blocked <n>` | Stop an affected step |
| `flowforge validate card <id>` | Verify the card |

## Turn Loop

1. Capture expected behavior, observed behavior, evidence, and affected FEATURE.
2. Classify using `classification-rules.md`.
3. Route a bug or accepted gap to an existing/new FEATURE; route reusable
   knowledge to a FIND or DEC. Never invent a legacy card route.
4. Record the result with `card log` and Proposal Journal.
5. Validate the affected card and leave unresolved questions visible.

For a design or compatibility gap, stop with `card steps <id> --status
blocked <n>`, record the blocker, and route it to `flowforge-design`. Do not
implement an unresolved product decision.

## Verification

Run `flowforge validate card <id>`, then the relevant tests and
`flowforge proposal inspect <proposal-id>`. Do not edit wiki or library files
directly, and do not claim a check passed without running it.
