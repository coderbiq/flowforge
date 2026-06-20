# FlowForge Agent Rules

此文件将在 `flowforge init` 时复制到目标项目根目录。

## CLI Conventions

- Use `--body -` with heredoc (`<<'EOF' ... EOF`) for multi-line card body content. Single-quoted heredoc delimiter prevents all shell expansion.
- Use `card create --batch <file>` for creating multiple cards at once.
- Use `card update --section "<name>" --body -` with heredoc to update a specific section.
- Use `-o json` to capture card IDs in scripts.
