# F-003 `.flowforge/state/` 是运行态和恢复态

- Status: validated
- Source: `docs/ARCHITECTURE.md`, `docs/GETTING-STARTED.md`

## Statement

FlowForge 的本地状态和恢复态存放在 `.flowforge/state/`，其职责是恢复活跃会话和工作流状态，而不是承载可再安装的工具本体。

## Why it matters

升级不应覆盖运行态数据，否则会打断工作恢复和会话连续性。

## References

- [Architecture](../../../docs/ARCHITECTURE.md)
- [Getting Started](../../../docs/GETTING-STARTED.md)
