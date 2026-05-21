# Configuration

`FlowForge` uses a machine-readable project config at `.flowforge/config.json`.

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

## Defaults

- `paths.tool_root`: `.flowforge`
- `paths.state_root`: `.flowforge/state`
- `task_backend.type`: `beads`
- `memory_provider.enabled`: `false`

## Notes

- Project identity is configuration, not code.
- Memory tags are derived from config and may be extended by adapters.
- Adapters may store additional runtime-only settings outside this file.
