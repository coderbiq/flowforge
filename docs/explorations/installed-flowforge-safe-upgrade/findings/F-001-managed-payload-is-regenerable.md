# F-001 安装产物的核心内容可重新生成

- Status: validated
- Source: `scripts/install.sh`, `workflow/templates/`, `workflow/guides/`, `workflow/schema/`

## Statement

安装后的 FlowForge 核心内容来自受管模板和脚本目录。`install.sh` 会把 workflow core、agents 和 scripts 安装到项目内的 `.flowforge/`，因此这些内容天然适合在升级时重新覆盖。

## Why it matters

如果安装产物本身就是模板化、可再生成的，那么升级策略就可以把它们视为工具自身的发布物，而不是项目业务数据。

## References

- [Installation guide](../../../GETTING-STARTED.md)
- [Installation script](../../../scripts/install.sh)
