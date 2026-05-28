# Configuration

`FlowForge` 使用一个机器可读的项目配置文件 `.flowforge/config.json`。

```json
{
  "project": {
    "id": "example-app",
    "name": "Example App",
    "slug": "example-app"
  },
  "paths": {
    "tool_root": ".flowforge",
    "state_root": ".flowforge/state"
  },
  "docs": {
    "default_workspace": "default",
    "workspaces": {
      "default": {
        "root": "docs",
        "scope": ".",
        "kind": "repository"
      }
    }
  },
  "task_backend": {
    "type": "beads"
  },
  "memory_provider": {
    "type": "memory-mcp",
    "enabled": true,
    "endpoint": "http://127.0.0.1:8000",
    "tags": ["project:example-app"]
  }
}
```

## 默认值

- `paths.tool_root`: `.flowforge`
- `paths.state_root`: `.flowforge/state`
- `task_backend.type`: `beads`
- `memory_provider.enabled`: `false`

## 约束

- 项目身份是 configuration，不是 code。
- memory tags 从 config 派生，也可以由 adapters 扩展。
- adapters 可以把额外的运行时设置放在这个文件之外。
- upgrade 操作必须保留这个文件，除非用户明确选择重新 bootstrap
  configuration。
- 如果要给 per-user memory 做认证，优先使用 `FLOWFORGE_MEMORY_ENDPOINT`
  和 `FLOWFORGE_MEMORY_API_KEY`。
- `OPENCODE_MEMORY_*` 仍然作为迁移期间的 legacy alias 被接受。
- 用户级 memory override 可以放在 `~/.config/flowforge/memory.json`，
  并按 project slug 或 id 做 key。
