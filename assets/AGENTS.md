<!-- FLOWFORGE:START -->
## FlowForge

Use `card init --type feature` to create cards; then edit the `.md` file directly for body content.
Use CLI for structured operations: `card link`, `card evolve`, `card log`, `card steps`.

### CLI
- `card init --type feature --title "..." --proposal <id>` to create a FEATURE card skeleton
- `card evolve <id> --stage designed|planned|done` for stage transitions (CLI enforces gates)
- `card log <id> --event "..." [--kind progress|bug|blocked]` to append to History
- `card steps <id> --status done|in_progress|blocked <n>` to update step status
- `context feature --feature <id> --step <n>` for minimal execution context
- `proposal inspect <id>` for auto-generated Feature Map and health checks
- `--body 'content\nwith\nnewlines'` for inline multi-line content
- Use single quotes for --body and --manifest to protect backticks, $, ! from shell expansion
- Never use shell redirects (`2>&1`, `<<`, `|`, `>`) with flowforge CLI — they trigger agent permission prompts
- `-o json` for machine-readable output
- `task`, `structure`, `log create` are DEPRECATED — use FEATURE-based commands instead

### Skills
| When | Skill |
|------|-------|
| Design / decompose proposal | `flowforge-design` |
| Execute implementation task | `flowforge-implement` |
| Report bug / finding / gap | `flowforge-feedback` |
| Import docs / archive proposal | `flowforge-curate` |
<!-- FLOWFORGE:END -->
