<!-- FLOWFORGE:START -->
## FlowForge

CLI is the only write path for cards. Never hand-write card files or frontmatter.

### CLI
- `--body 'content\nwith\nnewlines'` for multi-line content (\n for newlines, no shell redirects)
- Use single quotes for --body and --manifest to protect backticks, $, ! from shell expansion
- If body contains single quotes, break the string: `--body 'part1'"'"'part2'`
- `card batch --manifest 'cards:\n  - type: ...'` for multi-card creation (inline YAML with \n)
- Never use shell redirects (`2>&1`, `<<`, `|`, `>`) with flowforge CLI — they trigger agent permission prompts
- `-o json` for machine-readable output

### Skills
| When | Skill |
|------|-------|
| Design / decompose proposal | `flowforge-design` |
| Execute implementation task | `flowforge-implement` |
| Report bug / finding / gap | `flowforge-feedback` |
| Import docs / archive proposal | `flowforge-curate` |
<!-- FLOWFORGE:END -->
