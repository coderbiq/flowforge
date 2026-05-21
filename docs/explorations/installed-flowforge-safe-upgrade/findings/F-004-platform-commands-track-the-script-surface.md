# F-004 平台 commands 只是脚本 surface 的适配层

- Status: validated
- Source: `configs/claude/commands/flowforge/`, `configs/opencode/commands/flowforge/`, `workflow/guides/adapter-contract.md`

## Statement

Claude Code 和 OpenCode 的 `flowforge` commands 不是独立业务实现，而是围绕 `.flowforge/scripts/` 的平台包装层。它们的职责是把平台命令映射到同一套 FlowForge 脚本和工作流语义。

## Why it matters

如果脚本升级了，但平台 command 没有同步更新，用户在平台入口看到的行为、参数和说明会与实际工具版本脱节。

## References

- [Adapter contract](../../../workflow/guides/adapter-contract.md)
- `configs/claude/commands/flowforge/`
- `configs/opencode/commands/flowforge/`
