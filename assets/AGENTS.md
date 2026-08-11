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
- `journal append --actor <role> --message "..." [--references <card-id>] [--next "..."]` to record Proposal collaboration
- `journal recent [--proposal <id>] [--limit <n>]` to resume from recent collaboration notes
- `sync` detects OpenCode/Codex and reconciles Skills, subagents, routing rules, and the managed manifest
- `--body 'content\nwith\nnewlines'` for inline multi-line content
- Use single quotes for --body and --manifest to protect backticks, $, ! from shell expansion
- Never use shell redirects (`2>&1`, `<<`, `|`, `>`) with flowforge CLI — they trigger agent permission prompts
- `-o json` for machine-readable output
- `task`, `structure`, and `log create` are DEPRECATED; use FEATURE-based commands instead

### Skills
| When | Skill |
|------|-------|
| Design / decompose proposal | `flowforge-design` |
| Execute implementation task | `flowforge-implement` |
| Independently review a completed planned implementation | `flowforge-review` |
| Report bug / finding / gap | `flowforge-feedback` |
| Import docs / archive proposal | `flowforge-curate` |

### Subagent Orchestration

When FlowForge host subagents are installed, the Coordinator is a low-cost execution scheduler and the only interactive/delegating role. The Design Analyst owns framing, investigation planning, synthesis, and readiness decisions. The Investigator executes one registered brief and writes only its assigned FIND.

Before delegation, read structured Journal revision/readiness/re-entry state and tell the user what background action will run. Keep delegation one level deep: the Coordinator dispatches every worker directly, and workers never delegate or ask the user. External sources require explicit work-item authorization; unavailable required access returns `BLOCKED`.
<!-- FLOWFORGE:END -->
