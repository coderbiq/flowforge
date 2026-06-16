# flowforge-implement

Use when the user asks to execute a ready FlowForge implementation task, or provides a task id and wants code changes for that task. Do not use for design, analysis, feedback, archive, or general card lookup. Follow `references/workflow-rules.md`; only act on ready implementation tasks.

Card files are managed only through FlowForge CLI. Never edit wiki files, frontmatter, or links by hand.

## Hard Rules

- Always use single quotes for `--body` content containing mermaid, code blocks, or shell-special characters (`$`, `` ` ``, `!`, `{}`). Double-quoted `--body "..."` will be corrupted by shell expansion.
